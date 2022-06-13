import React, { useState, FC, useEffect } from 'react';
import { HStack } from '@chakra-ui/react';
import Button from 'components/ui/FormButton';
import FormLayout from 'layouts/FormLayout';
import NameIdentifier from '../NameIdentifier';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';

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
        heading={t`Name identifiers`}
        type={t`legal`}
        description={t`Enter the name and type of name by which the legal person is known. At least one legal name is required. Organizations are strongly encouraged to enter additional name identifiers such as Trading Name/ Doing Business As (DBA), Local names, and phonetics names where appropriate.`}
        ref={nameIdentifiersFieldArrayRef}
      />

      <NameIdentifier
        name="entity.name.local_name_identifiers"
        heading={t`Local Name Identifiers`}
        description={t`The name and type of name by which the legal person is known.`}
        ref={localNameIdentifiersFieldArrayRef}
      />

      <NameIdentifier
        name="entity.name.phonetic_name_identifiers"
        heading={t`Phonetic Name Identifiers`}
        description={t`The name and type of name by which the legal person is known.`}
        ref={phoneticNameIdentifiersFieldArrayRef}
      />

      <HStack width="100%" wrap="wrap" align="start" gap={4}>
        <Button borderRadius="5px" onClick={handleAddLegalNamesRow}>
          <Trans id="Add Legal Name">Add Legal Name</Trans>
        </Button>
        <Button borderRadius="5px" marginLeft="0 !important" onClick={handleAddNewLocalNamesRow}>
          <Trans id="Add Local Name">Add Local Name</Trans>
        </Button>
        <Button borderRadius="5px" marginLeft="0 !important" onClick={handleAddNewPhoneticNamesRow}>
          <Trans id="Add Phonetic Names">Add Phonetic Names</Trans>
        </Button>
      </HStack>
    </FormLayout>
  );
};

export default NameIdentifiers;
