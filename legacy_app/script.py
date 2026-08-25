import os
import json

class LegacyDataProcessor:
    def __init__(self, data_file):
        self.data_file = data_file

    def load_data(self):
        with open(self.data_file, 'r') as f:
            return json.load(f)

def run_processor():
    processor = LegacyDataProcessor("data.json")
    print(processor.load_data())
