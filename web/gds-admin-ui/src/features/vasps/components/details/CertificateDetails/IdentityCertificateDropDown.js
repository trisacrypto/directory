import PropTypes from 'prop-types';
import { Dropdown } from 'react-bootstrap';

const IdentityCertificateDropDown = ({ handleCopySignatureClick, handleCopySerialNumberClick }) => (
  <Dropdown className="float-end" align="end">
    <Dropdown.Toggle
      data-testid="certificate-details-3-dots"
      variant="link"
      tag="a"
      className="card-drop arrow-none cursor-pointer p-0 shadow-none"
    >
      <i className="dripicons-dots-3" />
    </Dropdown.Toggle>
    <Dropdown.Menu>
      <Dropdown.Item data-testid="copy-signature" onClick={handleCopySignatureClick}>
        <i className="mdi mdi-content-copy me-1" />
        Copy signature
      </Dropdown.Item>
      <Dropdown.Item data-testid="copy-serial-number" onClick={handleCopySerialNumberClick}>
        <i className="mdi mdi-content-copy me-1" />
        Copy serial number
      </Dropdown.Item>
    </Dropdown.Menu>
  </Dropdown>
);

IdentityCertificateDropDown.propTypes = {
  handleCopySignatureClick: PropTypes.func.isRequired,
  handleCopySerialNumberClick: PropTypes.func.isRequired,
};

export default IdentityCertificateDropDown;
