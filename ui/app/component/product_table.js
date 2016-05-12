import React from 'react';
import ObjectTable from 'cmpnt/table/object-table';

class ProductTable extends React.Component {
  filter(term, row) {
    let subMatch = row.sub_category &&
      row.sub_category.toLowerCase().indexOf(term.toLowerCase()) > -1;
    return (
      row.category.toLowerCase().indexOf(term.toLowerCase()) > -1 || subMatch
    );
  }

  render() {
    return (
      <div>
        <h2>Product Table</h2>
        <ObjectTable
          filter = {this.filter}
          rows = {this.props.rows}
          columns = {this.columns()}
          />
      </div>
    );
  }

  columns() {
    return [
      {field: 'id', label: 'add', format: id => {
        let disabled = this.props.productIds.filter(pid => pid === id)[0];
        return (
          <div>
            <button disabled={!this.props.hasStore} onClick={this.findPurchases.bind(this, id)}>
              Find Purchases
            </button>
            <button disabled={disabled} onClick={this.addProduct.bind(this, id)}>
              Add
            </button>
          </div>
        );
      }},
      {field: 'category', label: 'Category'},
      {field: 'sub_category', label: 'Sub Category'},
    ];
  }

  findPurchases(id) {
      this.props.onFindPurchases(id);
  }

  addProduct(id) {
    this.props.onProductAdd(id);
  }
}

ProductTable.propTypes = {
  rows: React.PropTypes.arrayOf(React.PropTypes.object).isRequired,
  hasStore: React.PropTypes.bool,
  onFindPurchases: React.PropTypes.func.isRequired,
  onProductAdd: React.PropTypes.func.isRequired,
  productIds: React.PropTypes.array.isRequired,
};

export default ProductTable
