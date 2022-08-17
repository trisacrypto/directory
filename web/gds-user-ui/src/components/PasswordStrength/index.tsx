import React, { useEffect, useState } from 'react';
import {
  Tr,
  Box,
  Text,
  Flex,
  Th,
  VStack,
  Stack,
  Icon,
  HStack,
  Heading,
  useColorModeValue
} from '@chakra-ui/react';

import { MdClose, MdDone } from 'react-icons/md';
import { Trans } from '@lingui/react';
const PasswordStrength = (props: any) => {
  const [isContains9Characters, setIsContains9Characters] = useState<boolean>(false);
  const [isContainsOneLowerCase, setIsContainsOneLowerCase] = useState<boolean>(false);
  const [isContainsOneUpperCase, setIsContainsOneUpperCase] = useState<boolean>(false);
  const [isContainsOneNumber, setIsContainsOneNumber] = useState<boolean>(false);
  const [isContainsOneSpecialChar, setIsContainsOneSpecialChar] = useState<boolean>(false);
  const checkPasswordValidity = (data: any) => {
    if (data.length >= 9) {
      setIsContains9Characters(true);
    } else {
      setIsContains9Characters(false);
    }
    // verify password contains at least one lowercase letter
    const lowerCaseLetters = /[a-z]/g;
    if (data.match(/^(?=.*[a-z]).*$/)) {
      setIsContainsOneLowerCase(true);
    } else {
      setIsContainsOneLowerCase(false);
    }

    // verify password contains at least one uppercase letter
    if (data.match(/^(?=.*[A-Z]).*$/)) {
      setIsContainsOneUpperCase(true);
    } else {
      setIsContainsOneUpperCase(false);
    }
    // verify password contains at least one number
    if (data.match(/[0-9]/)) {
      setIsContainsOneNumber(true);
    } else {
      setIsContainsOneNumber(false);
    }

    // verify password contains at least one special character
    const specialCharacters = /[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/;
    if (specialCharacters.test(data)) {
      setIsContainsOneSpecialChar(true);
    } else {
      setIsContainsOneSpecialChar(false);
    }
  };
  useEffect(() => {
    checkPasswordValidity(props.data);
  }, [props.data]);

  return (
    <Box>
      <Box>
        <Text
          textAlign={'left'}
          color={isContains9Characters ? 'gray.900' : 'gray.500'}
          data-testid="contains9Characters__text">
          <Text as={'span'} position={'relative'} top={2}>
            <Icon
              w={6}
              h={6}
              as={isContains9Characters ? MdDone : MdClose}
              color={isContains9Characters ? 'green' : 'gray.200'}
              data-testid="contains9Characters__icon"
            />{' '}
          </Text>
          <Trans id="At least 9 characters in length">At least 9 characters in length</Trans>
        </Text>
      </Box>
      <Box mt={2}>
        <Text fontWeight="semibold">
          <Trans id="Contain at least 3 of the following 4 types of characters">
            Contain at least 3 of the following 4 types of characters
          </Trans>
        </Text>
      </Box>
      <Box>
        <Text
          textAlign={'left'}
          color={isContainsOneLowerCase ? 'gray.900' : 'gray.500'}
          data-testid="containsOneLowerCase__text">
          <Text as={'span'} position={'relative'} top={2}>
            <Icon
              w={6}
              h={6}
              as={isContainsOneLowerCase ? MdDone : MdClose}
              color={isContainsOneLowerCase ? 'green' : 'gray.200'}
              data-testid="containsOneLowerCase__icon"
            />
          </Text>
          <Trans id="lower case letters (a-z)">lower case letters (a-z)</Trans>
        </Text>
      </Box>
      <Box>
        <Text
          textAlign={'left'}
          color={isContainsOneUpperCase ? 'gray.900' : 'gray.500'}
          data-testid="containsOneUpperCase__text">
          <Text as={'span'} position={'relative'} top={2}>
            <Icon
              w={6}
              h={6}
              as={isContainsOneUpperCase ? MdDone : MdClose}
              color={isContainsOneUpperCase ? 'green' : 'gray.200'}
              data-testid="containsOneUpperCase__icon"
            />{' '}
          </Text>
          <Trans id="upper case letters (A-Z)">upper case letters (A-Z)</Trans>
        </Text>
      </Box>
      <Box>
        <Text textAlign={'left'} color={isContainsOneNumber ? 'gray.900' : 'gray.500'}>
          <Text as={'span'} position={'relative'} top={2}>
            <Icon
              w={6}
              h={6}
              as={isContainsOneNumber ? MdDone : MdClose}
              color={isContainsOneNumber ? 'green' : 'gray.200'}
              data-testid="containsOneNumber__icon"
            />{' '}
          </Text>
          <Trans id="numbers (i.e. 0-9)">numbers (i.e. 0-9)</Trans>
        </Text>
      </Box>
      <Box>
        <Text textAlign={'left'} color={isContainsOneSpecialChar ? 'gray.900' : 'gray.500'}>
          <Text as={'span'} position={'relative'} top={2}>
            <Icon
              w={6}
              h={6}
              as={isContainsOneSpecialChar ? MdDone : MdClose}
              color={isContainsOneSpecialChar ? 'green' : 'gray.200'}
              data-testid="containsOneSpecialChar__icon"
            />{' '}
          </Text>
          <Trans id="special characters (e.g. !@#$%^&*)">special characters (e.g. !@#$%^&*)</Trans>
        </Text>
      </Box>
    </Box>
  );
};

export default PasswordStrength;
