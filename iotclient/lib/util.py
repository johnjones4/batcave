        
def make_std_headers(client_id: str, api_key: str):
    return {
        "X-Api-Key": api_key,
        "X-Client-Id": client_id
    }

