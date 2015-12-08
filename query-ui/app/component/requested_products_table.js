import React from 'react';
import ObjectTable from 'cmpnt/table/object-table';


class RequestedProductsTable extends React.Component {
  filter(term, row) {
    return row.name.indexOf(term) > -1
  }

  render() {
    let products = this.props.products
      .filter(product => this.props.requestedProducts[product.id])
      .map(product => {
        let name =
          `#{product.category}/#{product.sub_category?product.sub_category:""}`;
        let productCosts = this.props.requestedProducts[product.id];
        return {
          name,
          ...productCosts.stores,
        };
      });

    return(
      <ObjectTable
        filter = {this.filter}
        rows = {products}
        columns = {this.columns()}
      />
    );
  }

  columns() {
    firstField = {field: 'name', label: 'Name'};
    storeFields = this.props.stores.map(store => {
      return {
        field: `#{store.id}`,
        label: store.name,
        format: storeCosts => {
          if (!storeCosts) {
            return '--';
          }
          let result = [];
          for (let unit in storeCosts) {
            let {quantity, cost} = storeCosts[unit];
            result.push(`#{(cost/quantity).toFixed(2)} #{unit}`)
          }
          return result.join(', ');
        }
      };
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
