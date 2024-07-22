import React from 'react';
import { Button, HStack, Stack } from '@chakra-ui/react';
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
    <Stack data-testid="name-identifier">
      <NameIdentifier
        name="entity.name.name_identifiers"
        heading={t`Name Identifiers`}
        type="legal"
        description={t`Enter the name and type of name by which the legal person is known. At least one legal name is required. Organizations are strongly encouraged to enter additional name identifiers such as Trading Name/ Doing Business As (DBA), Local names, and phonetic names where appropriate.`}
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
        <Button onClick={handleAddLegalNamesRow}>
          <Trans id="Add Another Legal Name">Add Another Legal Name</Trans>
        </Button>
        <Button marginLeft="0 !important" onClick={handleAddNewLocalNamesRow}>
          <Trans id="Add Local Name">Add Local Name</Trans>
        </Button>
        <Button marginLeft="0 !important" onClick={handleAddNewPhoneticNamesRow}>
          <Trans id="Add Phonetic Names">Add Phonetic Names</Trans>
        </Button>
      </HStack>
    </Stack>
  );
};

export default NameIdentifiers;
