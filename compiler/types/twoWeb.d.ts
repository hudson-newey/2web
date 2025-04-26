// because $ is a compiler macro, we want to declare it so that TypeScript
// doesn't complain about it without adding it to the global runtime namespace
declare var $: void;
