import React from 'react';
import MultiSelect from 'cmpnt/select/multi';


let onRequestProductsClick = clickCallback => {
  return e => {
    e.stopPropagation();
    e.preventDefault();
    clickCallback();
  }
};

let mapProductToOption = {category, sub_category} => {
  let name = `${category}/${sub_category? sub_category+'/' : ''}`
  return {
    id: name,
    name,
  }
};

let ProductsToRequestInput = function({products, onProductsToRequestChange,
                                      onRequestProducts, productsToRequest}) {
  let options = products.map(mapProductToOption);
  let selected = productsToRequest.map(mapProductToOption);
  onRequestProducts = onRequestProductsClick(onRequestProducts);
  return(
    <div>
      <MultiSelect
        options={options}
        selected={selected}
        onChange={onRequestProducts}
      />
      <button
        type='submit'
        onClick={onRequestProducts}
      >
        Lookup
      </button>
    </div>
  );
};

OnProductsToRequestInput.propTypes = {
  products: React.PropTypes.arrayOf(React.PropTypes.shape({
    category:      React.PropTypes.string.isRequired,
    sub_category:  React.PropTypes.string,
  })).isRequired,
  productsToRequest: React.PropTypes.arrayOf(React.PropTypes.shape({
    category:      React.PropTypes.string.isRequired,
    sub_category:  React.PropTypes.string,
  })).isRequired,
  onProductsToRequestChange: React.PropTypes.func,
  onRequestProducts: React.PropTypes.func,
};


export default ProductsToRequestInput;
