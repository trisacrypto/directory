import { useEffect, useState } from "react";
import { copyToClipboard } from "../../utils";
import { Button } from "@chakra-ui/react";
import { Trans, t } from "@lingui/macro";

type CopyProps = {
    data: any;
};

function Copy({ data }: CopyProps) {
const memberDetailTableHeader = [
  {
    label: t`Name`,
    value: data.summary.name
  },
  {
    label: t`Website`,
    value: data.summary.website
  },
  {
    label: t`Business Category`,
    value: data.summary.business_category
  },
  {
    label: t`VASP Category`,
    value: data.summary.vasp_categories
  },
  {
    label: t`Country of Registration`,
    value: data.legal_person.country_of_registration
  },
  {
    label: t`Technical Contact`,
    value: data.contacts.technical
  },
  {
    label: t`Compliance / Legal Contact`,
    value: data.contacts.legal
  },
  {
    label: t`Administrative Contact`,
    value: data.contacts.administrative
  },
  {
    label: t`TRISA Endpoint`,
    value: data.summary.trisa_endpoint
  },
  {
    label: t`Common Name`,
    value: data.summary.common_name
  },
];
    const [copied, setCopied] = useState(false);

    useEffect(() => {
        const timeoutId = setTimeout(() => {
            setCopied(false);
        }, 2000);

        return () => {
            clearTimeout(timeoutId);
        };
    }
    , []);

    const handleCopy = async () => {
        // Remove label from the memberDetailTableHeader array
        await copyToClipboard(memberDetailTableHeader as any);
        setCopied(true);
    };
    return copied ? (
        <Button bg={'#FF7A59'} color={'white'}>
            <Trans>Copied</Trans>
        </Button>
    ) : (
        <Button bg={'#FF7A59'} color={'white'} onClick={handleCopy}>
            <Trans>Copy</Trans>
        </Button>
    );
}

export default Copy;
