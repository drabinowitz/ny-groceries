import React from 'react';
import request from '../ajax';

import ProductTable from './product_table';
import ProductsToRequestInput from './products_to_request_input';


export default class ProductCostsTable extends React.Component {
  constructor() {
    super()
    this.state = {
      products:           null,
      productsToRequest:  [],
      requestedProducts:  null,
      pendingRequest:     false,
    }
  }

  componentWillMount() {
    request('products/').then(products => this.setState({products}));
  }

  onRequestProducts() {
    this.setState({pendingRequest: true});
    new Promise((resolve, reject) => {
      let resolvedPromises = [];
      let rejected = false;
      this.state.productsToRequest.forEach((product, i, productsToRequest) => {
        request(`products/${product}/`)
          .then(data => {
            resolvedPromises.push(data.products);
            if (resolvedPromises.length === productsToRequest.length) {
              resolve(resolvedPromises);
            }
          })
          .catch(() => {
            if (!rejected) {
              rejected = true;
              reject();
            }
          });
      });
    }).then(productCosts => {
      let allProductCosts = productCosts.reduce((allCosts, costSet) => ({
        ...allCosts,
        ...costSet,
      }), this.state.requestedProducts || {});
      this.setState({
        pendingRequest: false,
        requestedProducts: allProductCosts,
        productsToRequest: [],
      });
    });
  }

  onProductsToRequestChange(productsToRequest) {
    this.setState({productsToRequest});
  }

  onClearAll(e) {
    e.stopPropagation();
    e.preventDefault();
    this.setState({requestedProducts: null});
  }

  render() {
    let overlayStyles = {
      "display":          this.state.pendingRequest? "block" :  "none",
      "position":         "fixed",
      "top":              0,
      "bottom":           0,
      "left":             0,
      "right":            0,
      "zIndex":           1000,
      "backgroundColor":  "rgba(100,100,100,0.5)"
    };

    if (!this.state.products) {
      return <h1>Loading App</h1>;
    }
    let requestedProductsTable = null;
    if (this.state.requestedProducts) {
      requestedProductsTable = <ProductTable
        rows={this.state.requestedProducts}
      />;
    }
    return(
      <div>
        <div className='overlay' style={overlayStyles}>
          <h1>Loading Product Costs</h1>
        </div>
        <div>
          <ProductsToRequestInput
            products                  = {this.state.products}
            productsToRequest         = {this.state.productsToRequest}
            onProductsToRequestChange =
              {this.onProductsToRequestChange.bind(this)}
            onRequestProducts         = {this.onRequestProducts.bind(this)}
          />
          <button
            onClick  = {this.onClearAll.bind(this)}
            disabled = {!this.state.requestedProducts}
          >
            Clear All
          </button>
          {requestedProductsTable}
        </div>
      </div>
    );
  }
}
