# crncl

A simple blog application written in Go.

## A note on styling
- utilizing Tailwind from CDN whilst in development; at deploy, we will want to include and minify and such

## TODO
- [ ] Rip out the custom stuff and push forward with a blog service.
    - [ ] image service feature
    - [ ] RSS feed support - https://kevincox.ca/2022/05/06/rss-feed-best-practices/
    - [ ] newsletter from posts
    - [ ] refactor templates for better componentizing and config passed in and coloring
        - [ ] work has started; need to componetize all templates and move to the template folder
    - [x] config for basic setup
    - [x] allow color theming from config (using tailwind colors and go-templating input)

## For personal when forking over
- [ ] Come up with color palette and apply it across - still needs applied to all page contents; nav is done
    - [ ] start with dark mode but also enable light or dark
- [ ] Add tags and filtering for blogs page
- [ ] Tooling to edit and add posts quickly
    - [ ] bonus to add the image generator to add an image
- [ ] Add a "subscribe" button for users to sign up to receive posts as email newsletter
- [ ] Where and how to deploy