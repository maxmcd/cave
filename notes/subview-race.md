

A page is loaded with two buttons, each in their own subcomponent.

When you click on the left button, it switches place with the right button.

If we simply reference subcomponents by their index then:

1. I double click on the left button
2. First event goes to the server server version swaps buttons
3. Second event goes to the server, it's on the first button, and might have a name from the first button, but it is sent to the second button's component because that's what's at the first index now.

So we could wait for a server response to do anything.
Or we could just let it be, a user can click on an incorrect version of the application and random things might happen :(

The alternative is still tough though right, let's say we let the double click happen and successfully send it to the right component? Well that would be dope. hmmm

oh, can I give them names!? yes, I can detect a move by looking at the memory address. ok so we use indexes and then... ugh this is complicated. Just indexes for now!


------

could count every component. get its memory address and put it in a hashmap. if we have seen a component before we don't increment the counter.

we pass that id back to the browser, the browser is associated with a component. if the component moves the id will follow it.
