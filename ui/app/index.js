import React from 'react';
import ReactDOM from 'react-dom';

import ReceiptForm from './component/receipt_form'

window.React = React;

var div = document.createElement('div');
div.id = 'app-container';
document.body.appendChild(div);

var style = document.createElement('link');
style.setAttribute('href', 'http://cmpnt.vistarmedia.com/app.css');
document.body.appendChild(style);

window.onload = function () {
  ReactDOM.render(<ReceiptForm />, div);
};
