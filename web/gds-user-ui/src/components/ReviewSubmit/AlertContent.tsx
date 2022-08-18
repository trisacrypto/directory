import { Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

const AlertContent = () => {
  return (
    <>
      <Text>
        <Text as={'span'}>
          <Trans id="Yes, I understand that this is the only time the PKCS12 Password is displayed and I have copied and securely saved the password.">
            Yes, I understand that this is the only time the PKCS12 Password is displayed and I have
            copied and securely saved the password.
          </Trans>
          <br />
          <Trans id="Click “No” if you have not copied the PKCS12 password yet and would like to view and copy the password.">
            Click “No” if you have not copied the PKCS12 password yet and would like to view and
            copy the password.
          </Trans>
          <br />
          <Trans id="Click “Yes” if you have copied the PKCS12 password and have securely saved it.">
            Click “Yes” if you have copied the PKCS12 password and have securely saved it.
          </Trans>
        </Text>
      </Text>
      <Text mt={4}>
        <Text as={'span'} fontWeight={'semibold'}>
          <Trans id="Note:">Note:</Trans>
        </Text>
        <Trans id="If you lose the PKCS12 password, you will have to the start the registration process from the beginning.">
          If you lose the PKCS12 password, you will have to the start the registration process from
          the beginning.
        </Trans>
      </Text>
    </>
  );
};

export default AlertContent;
