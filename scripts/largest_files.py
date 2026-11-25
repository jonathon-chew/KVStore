import os
import time
import queue
import threading
import argparse

from dataclasses import dataclass
from bisect import bisect_right


@dataclass
class Entry:
    file_path: str
    file_size: int


class LargestFiles:

    def __init__(self, size):
        self.size = size
        self.LOF: list[Entry] = []
        self.lock = threading.Lock()

    def _convert_bytes(self, number: int) -> str:
        if number == 0:
            return "0 b"

        for divisor, suffix in [(1000**3, "Gb"), (1000**2, "Mb"), (1000, "Kb")]:
            if number >= divisor:
                return f"{round(number / divisor, 2)} {suffix}"

        return f"{number} b"

    def add(self, entry: str, file_size: int):
        if self.size <= 0:
            return

        e = Entry(entry, file_size)

        with self.lock:
            if len(self.LOF) == self.size:
                if file_size < self.LOF[0].file_size:
                    return

            sizes = [item.file_size for item in self.LOF]
            position = bisect_right(sizes, file_size)

            self.LOF.insert(position, e)

            if len(self.LOF) > self.size:
                self.LOF.pop(0)

    def print_list(self):
        with self.lock:
            entries = list(reversed(self.LOF))

        for entry in entries:
            print(f"{entry.file_path}, " f"{self._convert_bytes(entry.file_size)}")

    def snapshot(self):
        with self.lock:
            return list(reversed(self.LOF))


class FileScanner:

    def __init__(self, track, workers=4):
        self.track = track
        self.workers = workers

        # Bounded queue prevents unlimited memory usage
        self.file_queue = queue.Queue(maxsize=1000)

        self.lock = threading.Lock()

        self.seen = 0
        self.processed = 0
        self.errors = 0

        self.finished = threading.Event()

    def discover(self):
        """Thread 1: Discover files and enqueue them."""

        try:
            for root, _, files in os.walk(".", onerror=self.walk_error):
                for file in files:
                    file_path = os.path.join(root, file)

                    with self.lock:
                        self.seen += 1

                    self.file_queue.put(file_path)

        finally:
            # Tell each worker there are no more files
            for _ in range(self.workers):
                self.file_queue.put(None)

    def walk_error(self, error):
        with self.lock:
            self.errors += 1

    def process(self):
        """Thread 2: Process files from the queue."""

        while True:
            file_path = self.file_queue.get()

            if file_path is None:
                break

            try:
                file_size = os.path.getsize(file_path)
                self.track.add(file_path, file_size)

            except OSError:
                with self.lock:
                    self.errors += 1

            finally:
                with self.lock:
                    self.processed += 1

    def display(self):
        """Thread 3: Display progress periodically."""

        start_time = time.time()

        while not self.finished.is_set():

            with self.lock:
                seen = self.seen
                processed = self.processed
                errors = self.errors

            elapsed = time.time() - start_time
            rate = processed / elapsed if elapsed else 0

            os.system("clear")

            print("=" * 50)
            print("         LARGEST FILES SCANNER")
            print("=" * 50)

            print(f"Files discovered : {seen:,}")
            print(f"Files processed  : {processed:,}")
            print(f"Files remaining  : {seen - processed:,}")
            print(f"Errors           : {errors:,}")
            print(f"Processing rate  : {rate:.2f} files/sec")
            print(f"Elapsed time     : {elapsed:.1f}s")

            print("=" * 50)

            time.sleep(1)

    def run(self):

        # Start discovery thread
        discovery = threading.Thread(target=self.discover)

        # Start worker threads
        workers = [threading.Thread(target=self.process) for _ in range(self.workers)]

        # Start display thread
        display = threading.Thread(target=self.display)

        display.start()

        for worker in workers:
            worker.start()

        discovery.start()

        # Wait for discovery to finish
        discovery.join()

        # Wait for all workers to finish
        for worker in workers:
            worker.join()

        # Stop display thread
        self.finished.set()
        display.join()

        print("\n\nScan complete!")
        print(f"Total discovered: {self.seen:,}")
        print(f"Total processed: {self.processed:,}")
        print(f"Total errors: {self.errors:,}")

        print("\nLargest files:")
        self.track.print_list()


def main():

    parser = argparse.ArgumentParser()
    parser.add_argument("--count", required=True, type=int)
    parser.add_argument("--workers", default=4, type=int)
    args = parser.parse_args()

    track = LargestFiles(args.count)
    scanner = FileScanner(track, workers=args.workers)

    scanner.run()


if __name__ == "__main__":
    main()
