export class AssetListing {
    assetId: string
    permissions: string[]

    constructor(init?: Partial<AssetListing>) {
        Object.assign(this, init)
    }
}
