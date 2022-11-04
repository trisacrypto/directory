import { CopyIcon } from '@chakra-ui/icons';
import {
  Box,
  Button,
  chakra,
  Flex,
  Heading,
  HStack,
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Stack,
  Text,
  VStack
} from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import useCertificateStepper from 'hooks/useCertificateStepper';
import FormLayout from 'layouts/FormLayout';
import { FiDownload } from 'react-icons/fi';
import { HiOutlineDotsHorizontal } from 'react-icons/hi';
import { useNavigate } from 'react-router-dom';

function CertificateInventory() {
  const navigate = useNavigate();
  const { jumpToLastStep } = useCertificateStepper();

  const handleEditClick = () => {
    navigate('/dashboard/certificate/registration');
    jumpToLastStep();
  };

  return (
    <>
      <Heading marginBottom="30px">
        <Trans>X.509 Identity Certificate Details</Trans>
      </Heading>
      <Stack spacing={5}>
        <HStack
          border="1px solid #DFE0EB"
          bg="#D8EAF6"
          py={3}
          px={8}
          borderRadius="10px"
          justifyContent="space-between">
          <Text fontWeight={700} fontSize="lg">
            <Trans>MAINNET Identity Certificate</Trans>
          </Text>
          <Button variant="primary" w="100%" maxW="130px" onClick={handleEditClick}>
            <Trans>Edit</Trans>
          </Button>
        </HStack>
        <FormLayout>
          <HStack justifyContent="space-between" w="100%">
            <Text textTransform="capitalize" fontWeight={700} fontSize="lg">
              <Trans>MainNet Certificate Details</Trans>
            </Text>
            <Menu>
              <MenuButton
                variant="ghost"
                color="#858585"
                as={Button}
                margin={0}
                marginInline={0}
                sx={{
                  '& .chakra-button__icon': {
                    marginInline: 0
                  }
                }}
                rightIcon={<HiOutlineDotsHorizontal size={25} />}
              />
              <MenuList>
                <MenuItem>
                  <CopyIcon mr={2} color="blue" /> Copy signature
                </MenuItem>
                <MenuItem>
                  <CopyIcon mr={2} color="blue" /> Copy serial number
                </MenuItem>
              </MenuList>
            </Menu>
          </HStack>
          <VStack align="start">
            <Text>
              <chakra.span fontWeight={700}>
                <Trans>Status</Trans>:
              </chakra.span>{' '}
              Verified
            </Text>
            <Text>
              <chakra.span fontWeight={700}>
                <Trans>Serial Number</Trans>:
              </chakra.span>{' '}
              S7NaVd8zt1YUEdwdfc7+Mg==
            </Text>
            <Text>
              <chakra.span fontWeight={700}>
                <Trans>Expires</Trans>:
              </chakra.span>{' '}
              Tue, 18 Apr 2023 21:14:39 GMT
            </Text>
            <Text>
              <chakra.span fontWeight={700}>
                <Trans>Issuer</Trans>:
              </chakra.span>{' '}
              CipherTrace Issuing CA
            </Text>
            <Text>
              <chakra.span fontWeight={700}>
                <Trans>Subject</Trans>:
              </chakra.span>{' '}
              trisa.alicevasp.io
            </Text>
            <Text>
              <chakra.span fontWeight={700}>
                <Trans>Endpoint</Trans>:
              </chakra.span>{' '}
              trisa.alicevasp.io:443
            </Text>
          </VStack>
        </FormLayout>
        <FormLayout>
          <Text textTransform="capitalize" fontWeight={700}>
            <Trans>Download</Trans>
          </Text>

          <Stack direction="row" justifyContent="start" w="100%" spacing={10}>
            <HStack
              border="1px solid #00000094"
              borderRadius="10px"
              p={3}
              alignItems="center!important"
              spacing={5}>
              <Flex gap={2}>
                <Flex
                  bg="#23A7E04D"
                  borderRadius="10px"
                  fontWeight={700}
                  color="rgba(85, 81, 81, 0.83)"
                  justifyContent="center"
                  alignItems="center"
                  p={2}>
                  .PEM
                </Flex>
                <Box>
                  <Text fontWeight={700}>Public Identity Key</Text>
                  <Text color="gray.600" fontSize="sm">
                    2.49kb
                  </Text>
                </Box>
              </Flex>
              <IconButton
                variant="ghost"
                fontSize="30px"
                color="blue"
                p={3}
                icon={<FiDownload />}
                aria-label="download"
              />
            </HStack>
            <HStack
              border="1px solid #00000094"
              borderRadius="10px"
              p={3}
              alignItems="center!important"
              spacing={5}>
              <Flex gap={2}>
                <Flex
                  bg="#23A7E04D"
                  borderRadius="10px"
                  fontWeight={700}
                  color="rgba(85, 81, 81, 0.83)"
                  justifyContent="center"
                  alignItems="center"
                  p={2}>
                  .GZ
                </Flex>
                <Box>
                  <Text fontWeight={700}>TRISA Trust Chain (CA)</Text>
                  <Text color="gray.600" fontSize="sm">
                    4.39 KB
                  </Text>
                </Box>
              </Flex>
              <IconButton
                variant="ghost"
                fontSize="30px"
                color="blue"
                p={3}
                icon={<FiDownload />}
                aria-label="download"
              />
            </HStack>
          </Stack>
        </FormLayout>
        <Stack direction="row" justifyContent="space-between" spacing={10}>
          <FormLayout w="100%">
            <HStack justifyContent="space-between" w="100%">
              <Text textTransform="capitalize" fontWeight={700} fontSize="lg">
                <Trans>User Details</Trans>
              </Text>
            </HStack>
            <VStack align="start">
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Common Name</Trans>:
                </chakra.span>{' '}
                Cyphertrace
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Country</Trans>:
                </chakra.span>{' '}
                United States
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Locality</Trans>:
                </chakra.span>{' '}
                Menlo Park
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Organization</Trans>:
                </chakra.span>{' '}
                CipherTrace
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Organization Unit</Trans>:
                </chakra.span>{' '}
                N/A
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Postal Code</Trans>:
                </chakra.span>{' '}
                N/A
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Province</Trans>:
                </chakra.span>{' '}
                California
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Serial Number</Trans>:
                </chakra.span>{' '}
                N/A
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Street Address</Trans>:
                </chakra.span>{' '}
                N/A
              </Text>
            </VStack>
          </FormLayout>
          <FormLayout w="100%">
            <HStack justifyContent="space-between" w="100%">
              <Text textTransform="capitalize" fontWeight={700} fontSize="lg">
                <Trans>Subject Details</Trans>
              </Text>
            </HStack>
            <VStack align="start">
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Common Name</Trans>:
                </chakra.span>{' '}
                Cyphertrace
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Country</Trans>:
                </chakra.span>{' '}
                United States
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Locality</Trans>:
                </chakra.span>{' '}
                Menlo Park
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Organization</Trans>:
                </chakra.span>{' '}
                CipherTrace
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Organization Unit</Trans>:
                </chakra.span>{' '}
                N/A
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Postal Code</Trans>:
                </chakra.span>{' '}
                N/A
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Province</Trans>:
                </chakra.span>{' '}
                California
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Serial Number</Trans>:
                </chakra.span>{' '}
                N/A
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Street Address</Trans>:
                </chakra.span>{' '}
                N/A
              </Text>
            </VStack>
          </FormLayout>
        </Stack>
      </Stack>
    </>
  );
}

export default CertificateInventory;
