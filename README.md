# shortinette
`shortinette` is a management system for programming assignments. It acts as a an orchestrator handling the complete lifecycle of assignments - from repository creation to automated grading.

`shortinette` is meant to be the open-source version of 42 Network's `Moulinette`, allowing anyone, 42 student or not, to organize their own learning events. 

`shortinette` automates assignment management by:
* Creating and configuring submission repositories on demand, including prefabricated development environments for students to focus on the learning rather than the painful configuration process of unfamiliar tech.
* Listening to repository webhook events to trigger automated gradings.
* Grading student submissions in safe, sandboxed environments.

We currently support an adapted version of Nils Mathieu's [Rust Piscine](https://github.com/nils-mathieu/piscine-rust), teaching students the basics of the Rust language, but `shortinette` is built to make the addition of new modules SO easy!

This is still a work in progress, as you read thisk, we are working on the technical documentation to allow you to seamlessly integrate your own assignments into the system. 

For more documentation on the process of organizing a Short (a learning event powered by `shortinette`), refer to this repo's [wiki](https://github.com/42-Short/shortinette/wiki).
