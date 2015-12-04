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
        if (!this.props.onProductAdd) {
          return id;
        }
        let disabled = this.props.productIds.filter(pid => pid === id)[0];
        return (
          <button disabled={disabled} onClick={this.addProduct.bind(this, id)}>
            Add
          </button>
        );
      }},
      {field: 'category', label: 'Category'},
      {field: 'sub_category', label: 'Sub Category'},
    ];
  }

  addProduct(id) {
    this.props.onProductAdd(id);
  }
}

ProductTable.propTypes = {
  rows: React.PropTypes.arrayOf(React.PropTypes.object).isRequired,
  onProductAdd: React.PropTypes.func,
  productIds: React.PropTypes.array,
};

export default ProductTable
