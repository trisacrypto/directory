import { FormControl, FormLabel, Select } from "@chakra-ui/react";
import { Trans } from "@lingui/macro";
import { useState } from "react";
import { DirectoryType, DirectoryTypeEnum } from "../memberType";
import { getMembersService } from "../service";

const MemberSelectNetwork = () => {
  const [selectNetwork, setSelectNetwork] = useState("");

  const handleSelectNetwork = (e: React.ChangeEvent<HTMLSelectElement>) => {
    e.preventDefault();
    return setSelectNetwork(e.target.value);
  };

  getMembersService(selectNetwork as DirectoryType);

    return (
    <FormControl>
      <FormLabel>
        <Trans>Select Network</Trans>
      </FormLabel>
      <Select
        data-testid="select-network"
        onChange={handleSelectNetwork}
      >
        <option value={DirectoryTypeEnum.MAINNET}>MainNet</option>
        <option value={DirectoryTypeEnum.TESTNET}>TestNet</option>
      </Select>
    </FormControl>
    );
};

export default MemberSelectNetwork;

