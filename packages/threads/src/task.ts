export type TaskCallback<Arguments extends unknown[], ReturnType> = (
  result: Arguments,
) => ReturnType;

export class Task<Arguments extends unknown[], ReturnType> {
  public constructor(
    private readonly callback: TaskCallback<Arguments, ReturnType>,
  ) {}

  public run(...args: Arguments): ReturnType {
    return this.callback(args);
  }
}
