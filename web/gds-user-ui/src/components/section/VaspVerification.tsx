import React from 'react';
import {
  Stack,
  Container,
  Box,
  Flex,
  Text,
  Link,
  chakra,
  FlexProps,
  StyleProps,
  useColorModeValue,
  UnorderedList,
  ListItem,
  Button
} from '@chakra-ui/react';

import { colors } from '../../utils/theme';

type Props = StyleProps &
  FlexProps & {
    children: React.ReactNode;
    title?: string;
  };

interface ILineProps {
  children?: React.ReactNode;
  title?: string;
}

const Line: React.FC<Props> = ({ children, title, ...rest }: any) => {
  return (
    <Flex>
      <Flex shrink={0}>
        <Flex rounded="md" bg={useColorModeValue('brand.500', '')} color="white"></Flex>
      </Flex>
      <Box ml={4}>
        <chakra.dt fontSize="lg" fontWeight="medium" lineHeight="6" {...rest}>
          {title}
        </chakra.dt>
        <chakra.dd mt={2}>{children}</chakra.dd>
      </Box>
    </Flex>
  );
};
export default function VaspVerification() {
  return (
    <Flex color={'black'} fontFamily={'Open Sans'} fontSize={'xl'} px={40}>
      <Stack>
        <Stack flex={1} justify={{ lg: 'center' }} py={{ base: 4, md: 20 }}>
          <Box my={{ base: 4 }} color="black">
            <Text fontFamily={'heading'} fontWeight={700} fontSize={'xl'}>
              TRISA members must complete a comprehensive multi-part verification form and due
              diligence process. Once verified, TRISA will issue TestNet and MainNet certificates
              for secure Travel Rule compliance.
            </Text>
          </Box>
          <Box bg={'gray.100'} p={5}>
            <Text fontSize={'xl'} color={'black'}>
              TRISA is the only global, open source, peer-to-peer and secure Travel Rule network.
              Become a TRISA-certified VASP today. Learn how TRISA works.
            </Text>
          </Box>
          <Box mt={20} pt={10}>
            <Stack
              display={{ md: 'grid' }}
              gridTemplateColumns={{ md: 'repeat(2,1fr)' }}
              color={'black'}
              gridColumnGap={{ md: 20, lg: 80 }}
              gridRowGap={{ md: 10 }}>
              <>
                <Line title="Sections & Details" fontWeight={'bold'}>
                  {''}
                </Line>
                <Line title="Who to Ask" fontWeight={'bold'}>
                  {''}
                </Line>
              </>

              <Line title="1 Basic Details" fontWeight={'bold'}>
                Information about the VASP such as website, incorporation date, business and VASP
                category.
              </Line>

              <Line>Business or Compliance Office</Line>

              <Line title="2 Legal Person" fontWeight={'bold'}>
                Information that identifies your organization as a Legal Person. This section
                represents the IVMS 101 data structure for legal persons and is strongly suggested
                for use as KYC information exchanged in TRISA transfers.
              </Line>
              <Line>Business or Compliance Office</Line>

              <Line title="3 Contacts" fontWeight={'bold'}>
                Contact information for representatives of your organization. Contacts include
                Technical, Legal/Compliance, Administrative, and Billing persons.
              </Line>
              <Line>Business or Compliance Office</Line>

              <Line title="4 TRISA Implementation" fontWeight={'bold'}>
                Technical information about your endpoint for certificate issuance. Each VASP is
                required to establish a TRISA endpoint for inter-VASP communication.
              </Line>
              <Line>Technical Officer</Line>
              <Line title="5 TRIXO Questionnaire" fontWeight={'bold'}>
                information to ensure that required compliance information exchanges are conducted
                correctly and safely. This includes information about jurisdiction and national
                regulator, CDD and Travel Rule policies, and data protection policies.
              </Line>
              <Line>Compliance Officer</Line>
              <Stack mt={20} bg={'gray.100'} py={5}>
                <Line title="Final Confirmation" fontWeight={'bold'}>
                  Upon submission, a member of TRISA’s verification team will review the form and
                  conduct a final due diligence phone call for physical verfication. Once due
                  diligence is complete, TRISA will issue certificates to the VASP.
                </Line>
              </Stack>
              <Stack mt={20} bg={'gray.100'} py={5}>
                <Line title="Need to Learn More?" fontWeight={'bold'}>
                  <UnorderedList>
                    <ListItem>
                      <Link>How TRISA Works</Link>
                    </ListItem>
                    <ListItem>
                      <Link>What is IVMS101?</Link>
                    </ListItem>
                  </UnorderedList>
                </Line>
              </Stack>
            </Stack>
          </Box>
          <Stack direction={['column', 'row']} pt={20} mx={10}>
            <Box>
              <Button
                bg={colors.system.blue}
                color={'white'}
                _hover={{
                  bg: '#10aaed'
                }}
                _focus={{
                  borderColor: 'transparent'
                }}>
                Download PDF
              </Button>
            </Box>
            <Box>
              <Button
                bg={colors.system.blue}
                color={'white'}
                _hover={{
                  bg: '#10aaed'
                }}
                _focus={{
                  borderColor: 'transparent'
                }}>
                Back to Getting Started
              </Button>
            </Box>
          </Stack>
        </Stack>
      </Stack>
    </Flex>
  );
}
