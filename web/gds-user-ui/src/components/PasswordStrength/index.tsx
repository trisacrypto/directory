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

import { MdClose } from 'react-icons/md';
const PasswordStrength = (props: any) => {
  console.log('watch password', props.data);
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
    if (data.match(/^(?=.*[~`!@#$%^&*()--+={}\[\]|\\:;"'<>,.?/_₹]).*$/)) {
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
        <Text textAlign={'left'} color={isContains9Characters ? 'gray.900' : 'gray.500'}>
          <Text as={'span'} position={'relative'} top={2}>
            <Icon w={6} h={6} as={MdClose} color={isContains9Characters ? 'green' : 'gray.200'} />{' '}
          </Text>
          At least 8 characters in length
        </Text>
      </Box>
      <Box mt={2}>
        <Text fontWeight="semibold">Contain at least 3 of the following 4 types of characters</Text>
      </Box>
      <Box>
        <Text textAlign={'left'} color={isContainsOneLowerCase ? 'gray.900' : 'gray.500'}>
          <Text as={'span'} position={'relative'} top={2}>
            <Icon w={6} h={6} as={MdClose} color={isContainsOneLowerCase ? 'green' : 'gray.200'} />{' '}
          </Text>
          lower case letters (a-z)
        </Text>
      </Box>
      <Box>
        <Text textAlign={'left'} color={isContainsOneLowerCase ? 'gray.900' : 'gray.500'}>
          <Text as={'span'} position={'relative'} top={2}>
            <Icon w={6} h={6} as={MdClose} color={isContainsOneUpperCase ? 'green' : 'gray.200'} />{' '}
          </Text>
          upper case letters (A-Z)
        </Text>
      </Box>
      <Box>
        <Text textAlign={'left'} color={isContainsOneNumber ? 'gray.900' : 'gray.500'}>
          <Text as={'span'} position={'relative'} top={2}>
            <Icon w={6} h={6} as={MdClose} color={isContainsOneNumber ? 'green' : 'gray.200'} />{' '}
          </Text>
          numbers (i.e. 0-9)
        </Text>
      </Box>
      <Box>
        <Text textAlign={'left'} color={isContainsOneSpecialChar ? 'gray.900' : 'gray.500'}>
          <Text as={'span'} position={'relative'} top={2}>
            <Icon
              w={6}
              h={6}
              as={MdClose}
              color={isContainsOneSpecialChar ? 'green' : 'gray.200'}
            />{' '}
          </Text>
          special characters (e.g. !@#$%^&*)
        </Text>
      </Box>
    </Box>
  );
};

export default PasswordStrength;
