export class AssetRequest {
    assetId: number;

    // For destroying
    managerAddr: string;

    // For freezing
    freezeAddr: string;
    freezeSetting: boolean;

    // For revoking
    clawbackAddr: string;
    recipientAddr: string;
    amount: number;

    // For freezing and revoking
    targetAddr: string;

    // For modifying
    currManagerAddr: string;
    newManagerAddr: string;
    newReserveAddr: string;
    newClawbackAddr: string;
    newFreezeAddr: string;

    constructor(init?: Partial<AssetRequest>) {
        Object.assign(this, init)
    }
}
