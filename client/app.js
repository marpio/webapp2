import React from "react"
import ReactDOM from "react-dom"
import ShoppingCart from "./shoppingcart/shoppingcart"

ReactDOM.render(<ShoppingCart url = {"/shoppingcart/items"} />,
  document.getElementById('shoppingcart-app')
);
