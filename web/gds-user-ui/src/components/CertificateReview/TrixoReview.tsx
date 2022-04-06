import React, { FC, useEffect } from 'react';
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
  Tag,
  TagLabel
} from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { getStepData } from 'utils/utils';
import { loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
interface TrixoReviewProps {}

const TrixoReview: React.FC<TrixoReviewProps> = (props) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [trixo, setTrixo] = React.useState<any>({});
  const getColorScheme = (status: string) => {
    if (status === 'yes' || status) {
      return 'cyan';
    } else {
      return '#eee';
    }
  };
  useEffect(() => {
    const getStepperData = loadDefaultValueFromLocalStorage();
    const stepData = {
      ...getStepperData.trixo
    };
    console.log('trixo step data', stepData);
    setTrixo(stepData);
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
          <Heading fontSize={24}>Section 5: TRIXO Questionnaire</Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            height={'34px'}
            onClick={() => jumpToStep(5)}
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
                <Td>Primary National Jurisdiction</Td>
                <Td>{trixo?.primary_national_jurisdiction}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Name of Primary Regulator</Td>
                <Td>{trixo?.primary_regulator}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Other Jurisdictions</Td>
                <Td>
                  {trixo?.other_jurisdictions?.map((o: any, i: any) => {
                    return (
                      <Text as={'span'} key={i}>
                        Country : {o.country} regulator name : {o.regulator_name}
                      </Text>
                    );
                  })}
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>CDD & Travel Rule Policies</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Does your organization have a programme that sets minimum AML, CFT, KYC/CDD and
                  Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes
                  where it is licensed/approved/registered?
                </Td>
                <Td>
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo.has_required_regulatory_program)}>
                    <TagLabel fontWeight={'bold'}>
                      {trixo.has_required_regulatory_program.toUpperCase()}
                    </TagLabel>
                  </Tag>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Does your organization conduct KYC/CDD before permitting its customers to
                  send/receive virtual asset transfers?
                </Td>
                <Td>
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo.financial_transfers_permitted)}>
                    <TagLabel fontWeight={'bold'}>
                      {trixo.financial_transfers_permitted.toUpperCase()}
                    </TagLabel>
                  </Tag>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Is your organization required to comply with the application of the Travel Rule
                  standards in the jurisdiction(s) where it is licensed/approved/registered?
                </Td>
                <Td>
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo.must_comply_travel_rule)}>
                    <TagLabel fontWeight={'bold'}>
                      {trixo.must_comply_travel_rule ? 'YES' : 'NO'}
                    </TagLabel>
                  </Tag>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>What is the minimum threshold for Travel Rule compliance?</Td>
                <Td>{`${trixo.kyc_threshold} ${trixo.kyc_threshold_currency}`}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Data Protection Policies</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Is your organization required by law to safeguard PII?</Td>
                <Td>
                  {' '}
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo.must_safeguard_pii)}>
                    <TagLabel fontWeight={'bold'}>
                      {trixo.must_safeguard_pii ? 'YES' : 'NO'}
                    </TagLabel>
                  </Tag>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Does your organization secure and protect PII, including PII received from other
                  VASPs under the Travel Rule?
                </Td>
                <Td>
                  {' '}
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo.safeguards_pii)}>
                    <TagLabel fontWeight={'bold'}>{trixo.safeguards_pii ? 'YES' : 'NO'}</TagLabel>
                  </Tag>
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
TrixoReview.defaultProps = {
  data: {}
};
export default TrixoReview;
