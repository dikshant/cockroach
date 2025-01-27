# Vector Indexing Demo

## How to use the files in this project

- Step 1: Download unplash lite embeddings: https://github.com/unsplash/datasets
- Step 2: run `python3 extract_embeddings.py`, make sure to place unsplash dataset in the same directory as the embedding extraction script 
- Step 3: run `python3 ./clipserver.py`
- Step 4: run `./dev build pkg/cmd/vecbench` and make sure embeddings generated in step 2 is in the bin/ directory
- Step 5: run `./bin/vecbench`

## Short method

- You can skip embedding extraction and use the embeddings in clip_embeddings.csv.zip by extracting it and placing it in the same directory as the vecbench binary
