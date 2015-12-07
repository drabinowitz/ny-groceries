import React from 'react';


class RequestedProductsTable extends React.Component {
  render() {
    let products = this.props.products
      .filter(product => this.props.requestedProducts[product.id])
      .sort((a, b) => {
        if (a.category != b.category) {
          return a.category > b.category;
        } else {
          return a.sub_category > b.sub_category;
        }})
      .map(product => {
        let name =
          `#{product.category}/#{product.sub_category?product.sub_category:""}`;
        let productCosts = this.props.requestedProducts[product.id];
      });
  }
}

RequestedProductsTable.propTypes = {
  requestedProducts:  React.PropTypes.objectOf(React.PropTypes.shape({
    product_id:       React.PropTypes.number.isRequired,
    stores:           React.PropTypes.objectOf(React.PropTypes.shape({
      store_id:       React.PropTypes.number.isRequired,
      units:          React.PropTypes.objectOf(React.PropTypes.shape({
        unit:         React.PropTypes.string.isRequired,
        quantity:     React.PropTypes.number.isRequired,
        cost:         React.PropTypes.number.isRequired,
      })).isRequired,
    })).isRequired,
  })).isRequired,

  products: React.PropTypes.arrayOf(React.PropTypes.shape({
    id:            React.PropTypes.number.isRequired,
    category:      React.PropTypes.string.isRequired,
    sub_category:  React.PropTypes.string,
  })).isRequired,

  stores: React.PropTypes.arrayOf(React.PropTypes.shape({
    id:    React.PropTypes.number.isRequired,
    name:  React.PropTypes.string.isRequired,
  })).isRequired,
}


export default RequestedProductsTable;
