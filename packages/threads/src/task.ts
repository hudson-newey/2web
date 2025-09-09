export type TaskCallback<Arguments extends any[], ReturnType> = (
  result: Arguments,
) => ReturnType;

export class Task<Arguments extends any[], ReturnType> {
  public constructor(
    private readonly callback: TaskCallback<Arguments, ReturnType>,
  ) {}

  public run(...args: Arguments): ReturnType {
    return this.callback(args);
  }
}
