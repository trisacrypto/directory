import { HStack, Button, FormControl } from '@chakra-ui/react';

import InputFormControl from 'components/ui/InputFormControl';
import { SubmitHandler, useForm } from 'react-hook-form';
import { t } from '@lingui/macro';
import { SearchIcon } from '@chakra-ui/icons';
import { Dispatch, SetStateAction, useEffect } from 'react';

type SearchVaspProps = {
  setSearchOrganization: Dispatch<SetStateAction<any>>;
};

const SearchVasp = ({ setSearchOrganization }: SearchVaspProps) => {
  // const { getAllOrganizations, organizations } = useOrganizationListByName();

  // const [showClose, setShowClose] = useState(false);
  const {
    handleSubmit,
    register,
    watch,
    formState: { errors }
  } = useForm();

  const onSubmit: SubmitHandler<any> = (data) => {
    setSearchOrganization(data.search);
  };

  useEffect(() => {
    setSearchOrganization(watch('search'));
  }, [watch('search')]);

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <FormControl>
        <HStack>
          <InputFormControl
            controlId="search"
            size="md"
            width={'100%'}
            isInvalid={!!errors.search}
            placeholder={t`VASP Name`}
            data-testid="name"
            type="search"
            maxLength={200}
            formHelperText={errors.search?.message as string}
            {...register('search')}
          />
          <Button variant="outline" type="submit" spinnerPlacement="start">
            <SearchIcon />
          </Button>
        </HStack>
      </FormControl>
    </form>
  );
};

export default SearchVasp;
