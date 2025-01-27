import torch
from torchvision import transforms
from PIL import Image
import clip
import pandas as pd
import os
from tqdm import tqdm
import numpy as np
import requests
from io import BytesIO
import logging
import time
import concurrent.futures
from functools import partial
import platform

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def setup_clip():
    """Initialize CLIP model and preprocessing optimized for M1"""
    try:
        logging.info("Setting up CLIP model...")
        
        # Check if running on Apple Silicon
        if platform.processor() == 'arm':
            device = torch.device("mps")
            logging.info("Using MPS (Metal Performance Shaders) device")
        else:
            device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
            logging.info(f"Using device: {device}")
            
        model, preprocess = clip.load("ViT-B/32", device=device)
        model = model.to(device)
        logging.info("CLIP model loaded successfully")
        
        return model, preprocess, device
    except Exception as e:
        logging.error(f"Error setting up CLIP: {str(e)}")
        # Fallback to CPU if MPS fails
        logging.info("Falling back to CPU")
        device = torch.device("cpu")
        model, preprocess = clip.load("ViT-B/32", device=device)
        return model, preprocess, device

def process_single_image(row, model, preprocess, device, session=None):
    """Process a single image with optimized memory handling"""
    url = row['photo_image_url']
    photo_id = row['photo_id']
    
    # Reuse session for better performance
    if session is None:
        session = requests.Session()
    
    try:
        response = session.get(url, timeout=10)
        response.raise_for_status()
        
        # Use BytesIO to avoid writing to disk
        image = Image.open(BytesIO(response.content)).convert('RGB')
        
        # Preprocess image
        image_input = preprocess(image).unsqueeze(0)
        
        # Move to appropriate device
        image_input = image_input.to(device)
        
        # Get embedding with automatic mixed precision
        with torch.no_grad():
            with torch.autocast(device_type=device.type if device.type != 'mps' else 'cpu'):
                image_features = model.encode_image(image_input)
                image_features /= image_features.norm(dim=-1, keepdim=True)
        
        # Move back to CPU and convert to numpy
        embedding = image_features.cpu().numpy()[0]
        
        return {
            'photo_id': photo_id,
            'photo_url': url,
            'description': row['photo_description'],
            'embedding': '[' + ','.join(map(str, embedding)) + ']',
            'success': True
        }
        
    except Exception as e:
        logging.error(f"Error processing {url}: {str(e)}")
        return {
            'photo_id': photo_id,
            'photo_url': url,
            'description': row['photo_description'],
            'embedding': None,
            'success': False
        }

def batch_processor(batch_df, model, preprocess, device):
    """Process a batch of images using a single session"""
    session = requests.Session()
    process_func = partial(process_single_image, model=model, preprocess=preprocess, 
                         device=device, session=session)
    return [process_func(row) for row in batch_df.to_dict('records')]

def save_batch(results, output_file, mode='a'):
    """Save batch results to CSV"""
    successful_results = [r for r in results if r['success']]
    if successful_results:
        df = pd.DataFrame(successful_results)
        df.to_csv(output_file, mode=mode, header=(mode=='w'), index=False)

def main():
    logging.info("Starting embedding extraction process")
    
    # Optimize for M1 Pro
    # M1 Pro has 8 performance cores + 2 efficiency cores
    # We'll use slightly fewer threads than cores to avoid overloading
    optimal_workers = 8
    batch_size = 32  # Smaller batches for better memory management
    
    model, preprocess, device = setup_clip()
    
    # Path to photos.tsv file
    photos_file = "unsplash-research-dataset-lite-latest/photos.tsv000"  # Update this
    output_file = "clip_embeddings.csv"
    
    # Read TSV file in chunks to manage memory
    chunk_size = 1000
    logging.info(f"Reading photos TSV file in chunks: {photos_file}")
    
    # Initialize output file
    pd.DataFrame(columns=['photo_id', 'photo_url', 'description', 'embedding'])\
        .to_csv(output_file, index=False)
    
    total_processed = 0
    total_failed = 0
    
    try:
        for chunk in pd.read_csv(photos_file, sep='\t', 
                               usecols=['photo_id', 'photo_image_url', 'photo_description'],
                               chunksize=chunk_size):
            
            logging.info(f"Processing chunk of {len(chunk)} rows")
            
            # Process in batches with multiple threads
            with concurrent.futures.ThreadPoolExecutor(max_workers=optimal_workers) as executor:
                futures = []
                
                for batch_start in range(0, len(chunk), batch_size):
                    batch_end = min(batch_start + batch_size, len(chunk))
                    batch_df = chunk.iloc[batch_start:batch_end]
                    
                    future = executor.submit(batch_processor, batch_df, model, preprocess, device)
                    futures.append(future)
                
                # Process completed batches
                for future in tqdm(concurrent.futures.as_completed(futures), 
                                 total=len(futures), 
                                 desc="Processing batches"):
                    results = future.result()
                    
                    # Count successes and failures
                    successful = sum(1 for r in results if r['success'])
                    failed = len(results) - successful
                    
                    total_processed += successful
                    total_failed += failed
                    
                    # Save batch results
                    save_batch(results, output_file)
            
            logging.info(f"Progress: processed={total_processed}, failed={total_failed}")
            
    except Exception as e:
        logging.error(f"Error processing file: {str(e)}")
    
    logging.info(f"Processing complete. Total processed: {total_processed}, Total failed: {total_failed}")

if __name__ == "__main__":
    main()