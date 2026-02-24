import type { Signal } from "../signal";

type ExecCallback<T> = () => T;
type ExecSubscriber = Signal<any>;

type ExecReturnType<T> = {
  returnValue: T;
  dependencies: Set<ExecSubscriber>;
};

const signalBufferNamespace = "__2_web_kit_signalBuffer";

/**
 * @description
 * Executes a callback and tracks all Signal subscriptions made during its
 * execution.
 */
export function execCallback<T>(callback: ExecCallback<T>): ExecReturnType<T> {
  const returnValue = callback();

  const signalBuffer = globalThis[signalBufferNamespace];

  return {
    returnValue: returnValue,
    dependencies: signalBuffer,
  };
}

export function addExecSubscriber(signal: ExecSubscriber) {
  globalThis[signalBufferNamespace].add(signal);
}

// This is a side effect to initialize the global signal buffer.
globalThis[signalBufferNamespace] ??= new Set<ExecSubscriber>();

declare global {
  var __2_web_kit_signalBuffer: Set<ExecSubscriber>;
}
