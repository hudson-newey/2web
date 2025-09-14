import worker from "./worker?url";

export class Thread extends Worker {
  public constructor() {
    super(worker);
  }
}
