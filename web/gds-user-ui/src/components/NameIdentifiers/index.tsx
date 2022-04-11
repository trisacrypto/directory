import React, { useState, FC, useEffect } from 'react';
import { HStack } from '@chakra-ui/react';
import Button from 'components/ui/FormButton';
import FormLayout from 'layouts/FormLayout';
import NameIdentifier from '../NameIdentifier';

const NameIdentifiers: React.FC = () => {
  const nameIdentifiersFieldArrayRef = React.useRef<any>(null);
  const localNameIdentifiersFieldArrayRef = React.useRef<any>(null);
  const phoneticNameIdentifiersFieldArrayRef = React.useRef<any>(null);

  const handleAddLegalNamesRow = () => {
    nameIdentifiersFieldArrayRef.current.addRow();
  };

  const handleAddNewLocalNamesRow = () => {
    localNameIdentifiersFieldArrayRef.current.addRow();
  };

  const handleAddNewPhoneticNamesRow = () => {
    phoneticNameIdentifiersFieldArrayRef.current.addRow();
  };

  return (
    <FormLayout>
      <NameIdentifier
        name="entity.name.name_identifiers"
        heading="Name identifiers"
        type={'legal'}
        description="The name and type of name by which the legal person is known."
        ref={nameIdentifiersFieldArrayRef}
      />

      <NameIdentifier
        name="entity.name.local_name_identifiers"
        heading="Local Name Identifiers"
        description="The name and type of name by which the legal person is known."
        ref={localNameIdentifiersFieldArrayRef}
      />

      <NameIdentifier
        name="entity.name.phonetic_name_identifiers"
        heading="Phonetic Name Identifiers"
        description="The name and type of name by which the legal person is known."
        ref={phoneticNameIdentifiersFieldArrayRef}
      />

      <HStack width="100%" wrap="wrap" align="start" gap={4}>
        <Button borderRadius="5px" onClick={handleAddLegalNamesRow}>
          Add Legal Name
        </Button>
        <Button borderRadius="5px" marginLeft="0 !important" onClick={handleAddNewLocalNamesRow}>
          Add Local Name
        </Button>
        <Button borderRadius="5px" marginLeft="0 !important" onClick={handleAddNewPhoneticNamesRow}>
          Add Phonetic Names
        </Button>
      </HStack>
    </FormLayout>
  );
};

export default NameIdentifiers;
