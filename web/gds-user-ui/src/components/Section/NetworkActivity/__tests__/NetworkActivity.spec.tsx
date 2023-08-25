import axios from "axios";
import { mockNetworkActivityData } from "../_mocks_";
import { networkActivity } from "../service";
import { render, screen } from "utils/test-utils";
import NetworkActivity from "../NetworkActivity";

function renderComponent() {
    return render(<NetworkActivity />);
}

describe('Network Activity Service', () => {
    it('returns network service data with response', async () => {
        const data = mockNetworkActivityData;
        axios.get = jest.fn().mockResolvedValue({ data });
        await expect (networkActivity()).resolves.toEqual(data);
    });
});

describe('Network Activity Component', () => {
    it('should render', () => {
        const { container } = renderComponent();
        expect(container).toMatchSnapshot();
    });

/*     it('should display the correct label on the y-axis', () => {
        renderComponent();
        const label = screen.getByDisplayValue('Network Activity');
        expect(label).toBeInTheDocument();
    });

    it('should display the testnet data', async () => {
        renderComponent();
        const { data } = mockNetworkActivityData;
       const testnetData = data?.networkActivity?.testnet
        expect(testnetData).toBeInTheDocument();
    });

    it('should display the mainnet data', async () => {
        renderComponent();
        const { data } = mockNetworkActivityData;
        const mainnetData = data?.networkActivity?.mainnet
        expect(mainnetData).toBeInTheDocument();
    }); */
});