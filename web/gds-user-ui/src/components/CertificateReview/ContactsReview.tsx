import React, { useEffect } from 'react';
import {
  Stack,
  Box,
  Text,
  Heading,
  Table,
  Tbody,
  Tr,
  Td,
  Button,
  useColorModeValue
} from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useSelector, RootStateOrAny } from 'react-redux';
import { hasValue } from 'utils/utils';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';

interface ContactsProps {
  data?: any;
}

const ContactsReview = (props: ContactsProps) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [contacts, setContacts] = React.useState<any>({});
  const textColor = useColorModeValue('gray.800', '#F7F8FC');

  useEffect(() => {
    const fetchData = async () => {
      const getStepperData = await getRegistrationDefaultValue();
      const stepData = {
        ...getStepperData.contacts
      };
      setContacts(stepData);
    };
    fetchData();
  }, [steps]);

  return (
    <Box
      border="2px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={textColor}
      fontSize={18}
      bg={useColorModeValue('white', '#171923')}
      p={5}
      px={5}>
      <Stack>
        <Box display={'flex'} justifyContent="space-between" pt={4} ml={0}>
          <Heading fontSize={20} mb="2rem">
            <Trans id="Section 3: Contacts">Section 3: Contacts</Trans>
          </Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            onClick={() => jumpToStep(3)}
            height={'34px'}
            _hover={{
              bg: '#10aaed'
            }}>
            <Trans id="Edit">Edit</Trans>
          </Button>
        </Box>
        <Stack fontSize={'1rem'}>
          <Table
            sx={{
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'semibold' },
              Tr: { borderStyle: 'hidden' },
              'td:first-child': {
                width: '50%',
                paddingLeft: '1rem'
              },
              td: {
                borderBottom: 'none',
                paddingInlineStart: 0,
                paddingY: 2.5
              }
            }}>
            <Tbody>
              {['technical', 'legal', 'administrative', 'billing'].map((contact, index) => (
                <Tr key={index}>
                  <Td textTransform="capitalize">
                    {t`${contact === 'legal' ? `Compliance / ${contact}` : contact} Contact`}
                  </Td>
                  <Td>
                    {hasValue(contacts?.[contact]) ? (
                      <>
                        {contacts?.[contact]?.name && <Text>{contacts?.[contact]?.name}</Text>}
                        {contacts?.[contact]?.email && <Text>{contacts?.[contact]?.email}</Text>}
                        {contacts?.[contact]?.phone && <Text>{contacts?.[contact]?.phone}</Text>}
                      </>
                    ) : (
                      'N/A'
                    )}
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Box>
  );
};
ContactsReview.defaultProps = {
  data: {}
};
export default ContactsReview;
