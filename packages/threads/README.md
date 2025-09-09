# 2Web Kit - Threads

True JavaScript multithreading support through workers.

## Tasks

Tasks are a high-level item that allows you to run a callback without having to
deal with threads.

```ts
const countToTen = new Task(() => {
  for (let i = 0; i < 11; i++) {
    console.log(i);
  }
});

await countToTen.run();
```

## Threads

```ts
const thread = new Thread();

const tasks = [
  new Task(() => { console.log("hello"); }),
  new Task(() => { console.log("world"); }),
];

for (const task of tasks) {
  thread.enqueue(tasks);
}
```

## Thread Pools

```ts
const threadPool = threadPool();

threadPool.workerCount;
threadPool.performanceMonitor;

threadPool.taskCount();

threadPool.start();
threadPool.stop();
threadPool.restart();
```
