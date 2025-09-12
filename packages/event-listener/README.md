# 2Web Kit - Event Listener

An event listener that automatically detaches itself when the associated element
is removed from the DOM.

Event listeners are a common cause of memory leaks because an attached event
listener will prevent an element node from being garbage collected.
