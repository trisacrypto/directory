import React, { FC, useEffect } from 'react';
import { Stack, Box, Text, Heading, Table, Tbody, Tr, Td, Button } from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { getStepData } from 'utils/utils';
import { loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
interface ContactsProps {
  data: any;
}

const ContactsReview = (props: ContactsProps) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [contacts, setContacts] = React.useState<any>({});
  useEffect(() => {
    const getStepperData = loadDefaultValueFromLocalStorage();
    const stepData = {
      ...getStepperData.contacts
    };
    console.log('contact step data', stepData);
    setContacts(stepData);
  }, [steps]);
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={'#252733'}
      bg={'white'}
      fontSize={18}
      p={5}
      px={5}>
      <Stack>
        <Box display={'flex'} justifyContent="space-between" pt={4} ml={5}>
          <Heading fontSize={24}>Section 3: Contacts</Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            onClick={() => jumpToStep(3)}
            height={'34px'}
            _hover={{
              bg: '#10aaed'
            }}>
            {' '}
            Edit{' '}
          </Button>
        </Box>
        <Stack fontSize={18}>
          <Table
            sx={{
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'bold' },
              Tr: { borderStyle: 'hidden' }
            }}>
            <Tbody>
              <Tr>
                <Td>Technical Contact</Td>
                <Td>
                  {contacts?.technical?.name} <br />
                  {contacts?.technical?.email} <br />
                  {contacts?.technical?.phone} <br />
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Compliance/ Legal Contact</Td>
                <Td>
                  {' '}
                  {contacts?.legal?.name} <br />
                  {contacts?.legal?.email} <br />
                  {contacts?.legal?.phone} <br />
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Administrative Contact</Td>
                <Td>
                  {' '}
                  {contacts?.administrative?.name} <br />
                  {contacts?.administrative?.email} <br />
                  {contacts?.administrative?.phone} <br />
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Billing Contact</Td>
                <Td>
                  {' '}
                  {contacts?.billing?.name} <br />
                  {contacts?.billing?.email} <br />
                  {contacts?.billing?.phone} <br />
                </Td>
                <Td></Td>
              </Tr>
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
