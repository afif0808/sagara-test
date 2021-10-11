<h1>Test : Sagara Backend Developer </h1>

<h3>Documentation</h3>
<h4>Steps to run the application</h4>
<ul>
  <li>Set the environtment variable in <code>.env</code> file accordingly </li>
  <li>If using Docker change <code>working_dir</code> and <code>volumes</code> part accordingly </li>
  <li>run the docker with <code>docker-compose</code> </li>
</ul>

<h4>Explanation</h4>

I'd like to give some notes and explanation about the application

<h5>Product's Image</h5>

For the product image,I made the endpoint for creating a new product and uploading the image separately.<br>
for the image storage location can be configured in <code>.env</code> file.

<h5>Authentication</h5>

For the authentication I made it stateless by using JWT, therefore users cannot log out since as far as i know in a pure stateless authentication , especially using JWT<br>
users cannot logout and there's no way to invalidate JWT token, I know we can do it with database or storage but then it's no longer completely stateless

<h5>Interface usage</h5>

As you might be already aware, that in this project I didn't define the interface in the same package as the implementation but instead in the package where it's used. <br>
I learned it recently that doing the former would mae the usage of interface less meaningful . 

https://www.efekarakus.com/golang/2019/12/29/working-with-interfaces-in-go.html













