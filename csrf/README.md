# Cross Site Request Forgery (CSRF)

# Design Specifications

## Title
Tricking users into performing unauthorized actions through Cross-Site Request Forgery (CSRF)

## Introduction
- Auction website that allows an attacker to perform a CSRF attack on another user. 

## Objectives
- prepare a legitimate site that allows a user to login and generate a session.
- display an auction page that has a vulnerable `<img>` tag that makes a forged request.
- perform an unintended action the legitimate site.

## Scope
- the auction site only has the minimum functionality handle a login.
- the vulnerable auction page can be hardcoded. 

## Components
- victims browser
- auction site
- forged request

## Mitigation Strategy
- Anti-CSRF tokens
- Headers
- Same-site Cookies

# Install
## Requirements
- `mysql`
- `php`
- `apache`

## MacOS
---
`sudo apachectl start`
`sudo mkdir /Library/WebServer/Documents/csrf`
`cp *.php /Library/WebServer/Documents/csrf`
open `http://localhost/csrf/index.php`
``



