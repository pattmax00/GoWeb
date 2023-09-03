# GoWeb 🌐

GoWeb is a simple Go web framework that aims to only use the standard library. The overall file structure and
development flow is inspired by larger frameworks like Laravel. It is partially ready for smaller projects if you are
fine with getting your hands dirty, but I plan on having it ready to go for more serious projects when it hits version
2.0.

<hr>

## Current features 🚀

- Routing/controllers
- Templating
- Simple database migration system
- Built in REST client
- CSRF protection
- Middleware
- Minimal user login/registration + sessions
- Config file handling
- Scheduled tasks
- Entire website compiles into a single binary (~10mb) (excluding env.json)
- Minimal dependencies (just standard library, postgres driver, and experimental package for bcrypt)

<hr>

## When to use 🙂

- You need to build a dynamic web application with persistent data
- You need to build a dynamic website using Go and need a good starting point
- You need to build an API in Go and don't know where to start
- Pretty much any use-case where you would use Laravel, Django, or Flask

## When not to use 🙃

- You need a static website (see [Hugo](https://gohugo.io/))
- You need a simple blog (see [Hugo](https://gohugo.io/))
- You need a simple site for your projects' documentation (see [Hugo](https://gohugo.io/))

## How to use 🤔

1. Clone
2. Delete the git folder, so you can start tracking in your own repo
3. Run `go get` to install dependencies
4. Copy env_example.json to env.json and fill in the values
5. Run `go run main.go` to start the server
6. Rename the occurences of "GoWeb" to your app name
7. Start building your app!
8. When you see useful changes to GoWeb you'd like in your project copy them over

## How to contribute 👨‍💻

- Open an issue on GitHub if you find a bug or have a feature request.
- [Email](mailto:contact@mpatterson.xyz) me a patch if you want to contribute code.
    - Please include a good description of what the patch does and why it is needed, also include how you want to be
      credited in the commit message.

<hr>

### License and disclaimer 😤

- You are free to use this project under the terms of the MIT license. See LICENSE for more details.
- You and you alone are responsible for the security and everything else regarding your application.
- It is not required, but I ask that when you use this project you give me credit by linking to this repository.
- I also ask that when releasing self-hosted or other end-user applications that you release it under
  the [GPLv3](https://www.gnu.org/licenses/gpl-3.0.html) license. This too is not required, but I would appreciate it.