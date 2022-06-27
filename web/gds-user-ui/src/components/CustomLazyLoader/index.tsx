export default function lazyLoaderPromise(lazyComponent: any, attemptsLeft = 3) {
  return new Promise((resolve, reject) => {
    lazyComponent()
      .then(resolve)
      .catch((error: any) => {
        // let us retry after 1500 ms
        setTimeout(() => {
          if (attemptsLeft === 1) {
            reject(error);
            return;
          }
          lazyLoaderPromise(lazyComponent, attemptsLeft - 1).then(resolve, reject);
        }, 1500);
      });
  });
}
