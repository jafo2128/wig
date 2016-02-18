## SSL

Because all of the information, including login details are passed over websockets we automatically generate SSL certificates to be used for the websocket connection on first start.

This can be a bit confusing and troublesome if you are not used to dealing with SSL certificates, but we feel that supporting unencrypted connections is just not an option, and as such we have tried to make the setup as straightforward as possible.

When you first start the wig server an SSL cert will be generated for localhost, this cert will be placed in the same directory as the executable and its path will be set up in the config file by default. 

If you wish to change the cert just replace these files with your cert and key and restart the server.

#### Valid Certificates

If you wish to move to a production setup and don't already have a valid cert for your domain, users of the site will see an SSL warning every time they visit the website, if you want to obtain a cert for free you can get one issued by Letsencrypt, otherwise you can purchase an SSL cert from any certificate authority and use that.

See [Letsencrypt](https://letsencrypt.org/) website for more info.

If you are using the ```letsencrypt-auto``` tool certs will be generated to the following location by default (on Debian).

```/etc/letsencrypt/live/{host}/fullchain.pem``` & ```/etc/letsencrypt/live/{host}/privkey.pem```

You can then simply update your config with these paths and your ready to rock! 