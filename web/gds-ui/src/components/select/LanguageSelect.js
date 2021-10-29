import React, { useContext } from 'react';
import Form from 'react-bootstrap/Form';
import LanguageContext from "../../contexts/LanguageContext";

const LanguageSelect = () => {
  const context = useContext(LanguageContext);

  return (
    <Form.Control
      as="select" custom
      value={context.language}
      onChange={e => context.changeLanguage(e.target.value)}
    >
      <option value="en">English</option>
      <option value="fr">French</option>
      <option value="de">German</option>
      <option value="zh">Chinese</option>
    </Form.Control>
  );
};

export default LanguageSelect;