export type AnimationIdentifier = () => symbol;

let animationIdCounter = 0;
export function animation(
  name = `@two-web/kit/animation (${animationIdCounter++})`,
): AnimationIdentifier {
  return () => Symbol(name);
}

