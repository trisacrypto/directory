import { Text } from '@chakra-ui/react';

const AlertContent = () => {
  return (
    <>
      <Text>
        <Text as={'span'}>
          Yes, I understand that this is the only time the PKCS12 Password is displayed and I have
          copied and securely saved the password. <br />
          Click “No” if you have not copied the PKCS12 password yet and would like to view and copy
          the password.
          <br />
          Click “Yes” if you have copied the PKCS12 password and have securely saved it.
        </Text>{' '}
      </Text>
      <Text mt={4}>
        <Text as={'span'} fontWeight={'semibold'}>
          Note:
        </Text>{' '}
        If you lose the PKCS12 password, you will have to the start the registration process from
        the beginning.
      </Text>
    </>
  );
};

export default AlertContent;
