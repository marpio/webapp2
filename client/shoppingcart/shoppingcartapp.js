import React from "react"
import ReactDOM from "react-dom"
import ShoppingCart from "./shoppingcart"

ReactDOM.render(<ShoppingCart url = {"/shoppingcart/items"} />,
  document.getElementById('shoppingcart-app')
);
