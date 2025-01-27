from flask import Flask, request, jsonify
import torch
import clip
import platform

app = Flask(__name__)

# Configure device for M1 Pro
if platform.processor() == 'arm':
    device = torch.device("mps")
    print("Using MPS (Metal Performance Shaders) device")
else:
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    print(f"Using device: {device}")

# Load CLIP model at startup
try:
    print("Loading CLIP model...")
    model, preprocess = clip.load("ViT-B/32", device=device)
    model = model.to(device)
    print("CLIP model loaded successfully")
except Exception as e:
    print(f"Error loading model: {e}")
    print("Falling back to CPU")
    device = torch.device("cpu")
    model, preprocess = clip.load("ViT-B/32", device=device)

@app.route('/embed-text', methods=['POST'])
def create_embedding():
    try:
        text = request.json.get('text')
        if not text:
            return jsonify({'error': 'No text provided'}), 400

        # Generate embedding
        with torch.no_grad():
            text_tokens = clip.tokenize([text]).to(device)
            with torch.autocast(device_type=device.type if device.type != 'mps' else 'cpu'):
                text_features = model.encode_text(text_tokens)
                text_features /= text_features.norm(dim=-1, keepdim=True)
            
        embedding = text_features.cpu().numpy()[0].tolist()
        return jsonify({'embedding': embedding})

    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/health', methods=['GET'])
def health_check():
    return jsonify({
        'status': 'healthy',
        'device': str(device)
    })

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)