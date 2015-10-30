import React from 'react';

class ProductForm extends React.Component {
  constructor() {
    super()
    this.state = {
      category: '',
      sub_category: '',
    };
  }

  render() {
    return (
      <div>
        <h4>Product Form</h4>
        <div>
          <input
            value={this.state.category}
            placeholder={'product category'}
            onChange={this.onCategoryChange.bind(this)}
            />
        </div>
        <div>
          <input
            value={this.state.sub_category}
            placeholder={'product sub category'}
            onChange={this.onSubCategoryChange.bind(this)}
            />
        </div>
        <button
          disabled={!this.state.category}
          onClick={this.onButtonClick.bind(this)}
          >
          Submit New Product
        </button>
      </div>
    );
  }

  onCategoryChange(e) {
    let category = e.target.value;
    e.preventDefault();
    e.stopPropagation();
    this.setState({category});
  }

  onSubCategoryChange(e) {
    let sub_category = e.target.value;
    e.preventDefault();
    e.stopPropagation();
    this.setState({sub_category});
  }

  onButtonClick() {
    this.props.onSubmit(this.state);
    this.setState({
      category: null,
      sub_category: null,
    });
  }
}

ProductForm.propTypes = {
  onSubmit: React.PropTypes.func.isRequired,
};

export default ProductForm;
