import {
  Avatar,
  Box,
  Button,
  Grid,
  GridItem,
  HStack,
  Select,
  Stack,
  StackDivider,
  Text
} from '@chakra-ui/react';
import { Account } from 'components/Account';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';

const OPTIONS = [
  { label: 'Newest registrations', value: 'NEWEST_REGISTRATIONS' },
  { label: 'Most recently logged in', value: 'MOST_RECENTLY_LOGGED_IN' },
  { label: 'Alphabetical', value: 'ALPHABETICAL' }
];

function ChooseAnAccount() {
  return (
    <Stack spacing={3} maxW="700px" mx="auto">
      <Grid templateColumns="repeat(5, 1fr)" gap={4}>
        <GridItem colSpan={3}>
          <InputFormControl controlId="search" />
        </GridItem>
        <GridItem colSpan={2}>
          <SelectFormControl controlId="filter" options={OPTIONS} />
        </GridItem>
      </Grid>
      <Stack divider={<StackDivider borderColor="#D9D9D9" />} p={2}>
        <Account username="John Doe" vaspName="VQSP Holding" />
        <Account username="John Doe" vaspName="VQSP Holding" />
        <Button textTransform="uppercase" fontWeight={700} borderRadius={0}>
          Load More
        </Button>
      </Stack>
    </Stack>
  );
}

export default ChooseAnAccount;