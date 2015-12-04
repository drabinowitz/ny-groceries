import React from 'react';
import MultiSelect from 'cmpnt/select/multi';


class ProductsToRequestInput extends React.Component {
  onRequestProductsClick(clickCallback) {
    return e => {
      e.stopPropagation();
      e.preventDefault();
      clickCallback();
    }
  }

  mapProductToOption({category, sub_category}) {
    let name = `${category}${sub_category? '/' + sub_category : ''}`
    return {
      id: name,
      name,
    }
  }

  render() {
    let options = this.props.products.map(this.mapProductToOption);
    let onRequestProducts = this.onRequestProductsClick(
      this.props.onRequestProducts);
    return(
      <div>
        <MultiSelect
          options  = {options}
          value    = {this.props.productsToRequest}
          onChange = {this.props.onProductsToRequestChange}
        />
        <button
          type    = 'submit'
          onClick = {onRequestProducts}
        >
          Lookup
        </button>
      </div>
    );
  }
}

ProductsToRequestInput.propTypes = {
  products: React.PropTypes.arrayOf(React.PropTypes.shape({
    category:      React.PropTypes.string.isRequired,
    sub_category:  React.PropTypes.string,
  })).isRequired,
  productsToRequest: React.PropTypes.arrayOf(React.PropTypes.string).isRequired,
  onProductsToRequestChange: React.PropTypes.func,
  onRequestProducts: React.PropTypes.func,
};


export default ProductsToRequestInput;
