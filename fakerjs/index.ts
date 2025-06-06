import { faker } from "@faker-js/faker";
import {
  CommentType,
  PrismaClient,
  Gender,
  Brand,
  Account,
  ProductModel,
  Product,
  PaymentMethod,
  Status,
  RefundMethod,
  Prisma,
  Tag,
  ProductType,
  Refund,
  ResourceType,
  AdminAccount,
  UserAccount,
  Address,
  Payment,
  SaleType,
  Sale,
  AccountType,
} from "../node_modules/.prisma/client";

type TxPrisma = Omit<
  PrismaClient,
  "$connect" | "$disconnect" | "$on" | "$transaction" | "$use" | "$extends"
>;

// Permission enum
enum Permission {
  // Product permissions
  PERMISSION_CREATE_PRODUCT = "PERMISSION_CREATE_PRODUCT",
  PERMISSION_UPDATE_PRODUCT = "PERMISSION_UPDATE_PRODUCT",
  PERMISSION_DELETE_PRODUCT = "PERMISSION_DELETE_PRODUCT",
  PERMISSION_VIEW_PRODUCT = "PERMISSION_VIEW_PRODUCT",

  // Product Model permissions
  PERMISSION_CREATE_PRODUCT_MODEL = "PERMISSION_CREATE_PRODUCT_MODEL",
  PERMISSION_UPDATE_PRODUCT_MODEL = "PERMISSION_UPDATE_PRODUCT_MODEL",
  PERMISSION_DELETE_PRODUCT_MODEL = "PERMISSION_DELETE_PRODUCT_MODEL",
  PERMISSION_VIEW_PRODUCT_MODEL = "PERMISSION_VIEW_PRODUCT_MODEL",

  // Product Serial permissions
  PERMISSION_CREATE_PRODUCT_SERIAL = "PERMISSION_CREATE_PRODUCT_SERIAL",
  PERMISSION_UPDATE_PRODUCT_SERIAL = "PERMISSION_UPDATE_PRODUCT_SERIAL",
  PERMISSION_DELETE_PRODUCT_SERIAL = "PERMISSION_DELETE_PRODUCT_SERIAL",
  PERMISSION_VIEW_PRODUCT_SERIAL = "PERMISSION_VIEW_PRODUCT_SERIAL",

  // Sale permissions
  PERMISSION_CREATE_SALE = "PERMISSION_CREATE_SALE",
  PERMISSION_UPDATE_SALE = "PERMISSION_UPDATE_SALE",
  PERMISSION_DELETE_SALE = "PERMISSION_DELETE_SALE",
  PERMISSION_VIEW_SALE = "PERMISSION_VIEW_SALE",

  // Tag permissions
  PERMISSION_CREATE_TAG = "PERMISSION_CREATE_TAG",
  PERMISSION_UPDATE_TAG = "PERMISSION_UPDATE_TAG",
  PERMISSION_DELETE_TAG = "PERMISSION_DELETE_TAG",
  PERMISSION_VIEW_TAG = "PERMISSION_VIEW_TAG",

  // Brand permissions
  PERMISSION_CREATE_BRAND = "PERMISSION_CREATE_BRAND",
  PERMISSION_UPDATE_BRAND = "PERMISSION_UPDATE_BRAND",
  PERMISSION_DELETE_BRAND = "PERMISSION_DELETE_BRAND",
  PERMISSION_VIEW_BRAND = "PERMISSION_VIEW_BRAND",

  // Payment permissions
  PERMISSION_UPDATE_PAYMENT = "PERMISSION_UPDATE_PAYMENT",
  PERMISSION_DELETE_PAYMENT = "PERMISSION_DELETE_PAYMENT",
  PERMISSION_VIEW_PAYMENT = "PERMISSION_VIEW_PAYMENT",

  // Refund permissions
  PERMISSION_UPDATE_REFUND = "PERMISSION_UPDATE_REFUND",
  PERMISSION_DELETE_REFUND = "PERMISSION_DELETE_REFUND",
  PERMISSION_VIEW_REFUND = "PERMISSION_VIEW_REFUND",
  PERMISSION_APPROVE_REFUND = "PERMISSION_APPROVE_REFUND",

  // Comment permissions
  PERMISSION_CREATE_COMMENT = "PERMISSION_CREATE_COMMENT",
  PERMISSION_UPDATE_COMMENT = "PERMISSION_UPDATE_COMMENT",
  PERMISSION_DELETE_COMMENT = "PERMISSION_DELETE_COMMENT",
  PERMISSION_VIEW_COMMENT = "PERMISSION_VIEW_COMMENT",

  // User permissions
  PERMISSION_VIEW_USER = "PERMISSION_VIEW_USER",
  PERMISSION_UPDATE_USER = "PERMISSION_UPDATE_USER",
  PERMISSION_DELETE_USER = "PERMISSION_DELETE_USER",

  // Address permissions
  PERMISSION_VIEW_ADDRESS = "PERMISSION_VIEW_ADDRESS",
  PERMISSION_UPDATE_ADDRESS = "PERMISSION_UPDATE_ADDRESS",
  PERMISSION_DELETE_ADDRESS = "PERMISSION_DELETE_ADDRESS",
}

// Role and Permission definitions
const ROLES = {
  PRODUCT_MANAGER: {
    id: "PRODUCT_MANAGER",
    description: "Can manage products, models, and brands",
    permissions: [
      Permission.PERMISSION_CREATE_PRODUCT,
      Permission.PERMISSION_UPDATE_PRODUCT,
      Permission.PERMISSION_DELETE_PRODUCT,
      Permission.PERMISSION_VIEW_PRODUCT,
      Permission.PERMISSION_CREATE_PRODUCT_MODEL,
      Permission.PERMISSION_UPDATE_PRODUCT_MODEL,
      Permission.PERMISSION_DELETE_PRODUCT_MODEL,
      Permission.PERMISSION_VIEW_PRODUCT_MODEL,
      Permission.PERMISSION_CREATE_BRAND,
      Permission.PERMISSION_UPDATE_BRAND,
      Permission.PERMISSION_DELETE_BRAND,
      Permission.PERMISSION_VIEW_BRAND,
    ],
  },
  INVENTORY_MANAGER: {
    id: "INVENTORY_MANAGER",
    description: "Can manage inventory and product serials",
    permissions: [
      Permission.PERMISSION_CREATE_PRODUCT_SERIAL,
      Permission.PERMISSION_UPDATE_PRODUCT_SERIAL,
      Permission.PERMISSION_DELETE_PRODUCT_SERIAL,
      Permission.PERMISSION_VIEW_PRODUCT_SERIAL,
      Permission.PERMISSION_VIEW_PRODUCT,
      Permission.PERMISSION_VIEW_PRODUCT_MODEL,
    ],
  },
  SALE_MANAGER: {
    id: "SALE_MANAGER",
    description: "Can manage sales and discounts",
    permissions: [
      Permission.PERMISSION_CREATE_SALE,
      Permission.PERMISSION_UPDATE_SALE,
      Permission.PERMISSION_DELETE_SALE,
      Permission.PERMISSION_VIEW_SALE,
      Permission.PERMISSION_VIEW_PRODUCT,
      Permission.PERMISSION_VIEW_PRODUCT_MODEL,
    ],
  },
  PAYMENT_MANAGER: {
    id: "PAYMENT_MANAGER",
    description: "Can manage payments and refunds",
    permissions: [
      Permission.PERMISSION_UPDATE_PAYMENT,
      Permission.PERMISSION_DELETE_PAYMENT,
      Permission.PERMISSION_VIEW_PAYMENT,
      Permission.PERMISSION_UPDATE_REFUND,
      Permission.PERMISSION_DELETE_REFUND,
      Permission.PERMISSION_VIEW_REFUND,
      Permission.PERMISSION_APPROVE_REFUND,
    ],
  },
  CONTENT_MANAGER: {
    id: "CONTENT_MANAGER",
    description: "Can manage content (tags, comments)",
    permissions: [
      Permission.PERMISSION_CREATE_TAG,
      Permission.PERMISSION_UPDATE_TAG,
      Permission.PERMISSION_DELETE_TAG,
      Permission.PERMISSION_VIEW_TAG,
      Permission.PERMISSION_CREATE_COMMENT,
      Permission.PERMISSION_UPDATE_COMMENT,
      Permission.PERMISSION_DELETE_COMMENT,
      Permission.PERMISSION_VIEW_COMMENT,
    ],
  },
  USER_MANAGER: {
    id: "USER_MANAGER",
    description: "Can manage user accounts and addresses",
    permissions: [
      Permission.PERMISSION_VIEW_USER,
      Permission.PERMISSION_UPDATE_USER,
      Permission.PERMISSION_DELETE_USER,
      Permission.PERMISSION_VIEW_ADDRESS,
      Permission.PERMISSION_UPDATE_ADDRESS,
      Permission.PERMISSION_DELETE_ADDRESS,
    ],
  },
} as const;

async function createRoles(prisma: TxPrisma) {
  // Create roles
  for (const role of Object.values(ROLES)) {
    await prisma.role.create({
      data: {
        id: role.id,
        description: role.description,
      },
    });

    // Create permissions for each role
    for (const permissionId of role.permissions) {
      // Create permission if it doesn't exist
      await prisma.permission.upsert({
        where: { id: permissionId },
        create: {
          id: permissionId,
          description: `Permission to ${permissionId
            .replace(/PERMISSION_/g, "")
            .replace(/_/g, " ")
            .toLowerCase()}`,
        },
        update: {},
      });

      // Create role-permission relationship
      await prisma.permissionOnRole.create({
        data: {
          role_id: role.id,
          permission_id: permissionId,
        },
      });
    }
  }
}

async function createAccounts(prisma: TxPrisma, count: number) {
  const adminAccounts: AdminAccount[] = [];
  const userAccounts: UserAccount[] = [];

  // Create admin accounts
  for (let i = 0; i < Math.ceil(count * 0.1); i++) {
    const account = await prisma.account.create({
      data: {
        username: `admin${i + 1}`,
        password: "admin123", // In production, this should be hashed
        type: AccountType.ADMIN,
        admin_account: {
          create: {
            is_super_admin: i === 0,
            avatar_url: faker.image.avatar(),
          },
        },
      },
      include: {
        admin_account: true,
      },
    });

    // Assign random roles to admin
    const roleIds = Object.keys(ROLES);
    const numRoles = Math.floor(Math.random() * 3) + 1; // 1-3 roles per admin
    const selectedRoles = faker.helpers.arrayElements(roleIds, numRoles);

    for (const roleId of selectedRoles) {
      await prisma.roleOnAdmin.create({
        data: {
          admin_id: account.id,
          role_id: roleId,
        },
      });
    }

    adminAccounts.push(account.admin_account!);
  }

  // Create user accounts
  for (let i = 0; i < Math.floor(count * 0.9); i++) {
    const account = await prisma.account.create({
      data: {
        username: `user${i + 1}`,
        password: "user123", // In production, this should be hashed
        type: AccountType.USER,
        user_account: {
          create: {
            email: faker.internet.email(),
            phone: faker.phone.number(),
            gender: faker.helpers.enumValue(Gender),
            full_name: faker.person.fullName(),
            avatar_url: faker.image.avatar(),
          },
        },
      },
      include: {
        user_account: true,
      },
    });

    userAccounts.push(account.user_account!);
  }

  return { adminAccounts, userAccounts };
}

// Generate Brands
async function createBrands(prisma: TxPrisma, count: number) {
  const brandsData = Array.from({ length: count }, () => ({
    name: faker.company.name(),
    description: faker.company.catchPhrase(),
  }));

  await prisma.brand.createMany({
    data: brandsData,
    skipDuplicates: true,
  });
  return await prisma.brand.findMany({ take: count });
}

// Generate Tags
async function createTags(prisma: TxPrisma, count: number) {
  const tagsData = Array.from({ length: count }, () => ({
    tag: faker.commerce.department(),
    description: faker.commerce.productDescription(),
  }));

  await prisma.tag.createMany({
    data: tagsData,
    skipDuplicates: true,
  });
  return await prisma.tag.findMany({ take: count });
}

// Generate ProductTypes
async function createProductTypes(prisma: TxPrisma, count: number) {
  const productTypesData = Array.from({ length: count }, () => ({
    name: faker.commerce.productMaterial(),
  }));

  await prisma.productType.createMany({
    data: productTypesData,
    skipDuplicates: true,
  });
  return await prisma.productType.findMany({ take: count });
}

// Generate ProductModels
async function createProductModels(
  prisma: TxPrisma,
  brands: Brand[],
  tags: Tag[],
  productTypes: ProductType[],
  count: number
) {
  const productModelsData: Prisma.ProductModelCreateManyInput[] = [];
  const tagOnProductData: Prisma.TagOnProductModelCreateManyInput[] = [];

  for (let i = 0; i < count; i++) {
    // Randomly assign 1-3 tags to each product model
    const tagCount = Math.floor(Math.random() * 3) + 1;
    const selectedTags = faker.helpers.arrayElements(tags, tagCount);
    const modelName = faker.commerce.productName();

    productModelsData.push({
      brand_id: brands[Math.floor(Math.random() * brands.length)].id,
      type: productTypes[Math.floor(Math.random() * productTypes.length)].id,
      name: modelName,
      description: faker.commerce.productDescription(),
      list_price: BigInt(
        parseInt(faker.commerce.price({ min: 100, max: 1000 })) * 1000
      ),
      date_manufactured: faker.date.past(),
    });

    // We'll need to fetch the created product models to get their IDs for the tag connections
  }

  await prisma.productModel.createMany({
    data: productModelsData,
    skipDuplicates: true,
  });

  // Fetch created product models to get their IDs
  const createdProductModels = await prisma.productModel.findMany({
    take: count,
  });

  // Create tag connections
  for (const productModel of createdProductModels) {
    const tagCount = Math.floor(Math.random() * 3) + 1;
    const selectedTags = faker.helpers.arrayElements(tags, tagCount);

    for (const tag of selectedTags) {
      tagOnProductData.push({
        product_model_id: productModel.id,
        tag: tag.tag,
      });
    }
  }

  if (tagOnProductData.length > 0) {
    await prisma.tagOnProductModel.createMany({
      data: tagOnProductData,
      skipDuplicates: true,
    });
  }

  return createdProductModels;
}

// Generate Products
async function createProducts(
  prisma: TxPrisma,
  productModels: ProductModel[],
  count: number
) {
  const products: Product[] = [];

  for (let i = 0; i < count; i++) {
    const product = await prisma.product.create({
      data: {
        product_model_id:
          productModels[Math.floor(Math.random() * productModels.length)].id,
        additional_price: BigInt(Math.floor(Math.random() * 50000)),
        is_active: Math.random() < 0.9, // 90% chance of being active
        can_combine: Math.random() < 0.3, // 30% chance of being combinable
        metadata: {
          color: faker.color.human(),
          size: Math.floor(Math.random() * 5) + 1,
          weight: Math.floor(Math.random() * 1000) + 100, // 100-1100g
          dimensions: {
            length: Math.floor(Math.random() * 50) + 10, // 10-60cm
            width: Math.floor(Math.random() * 30) + 5, // 5-35cm
            height: Math.floor(Math.random() * 20) + 5, // 5-25cm
          },
        },
        ProductTracking: {
          create: {
            current_stock: BigInt(Math.floor(Math.random() * 100) + 10),
            sold: BigInt(Math.floor(Math.random() * 5)),
          },
        },
      },
      include: {
        ProductTracking: true,
      },
    });

    // Create product serials for each product
    const serialsData = Array.from(
      { length: Number(product.ProductTracking?.current_stock ?? 0) },
      () => ({
        serial_id: faker.string.alphanumeric(10).toUpperCase(),
        product_id: product.id,
        is_sold: false,
        is_active: true,
      })
    );

    await prisma.productSerial.createMany({
      data: serialsData,
      skipDuplicates: true,
    });

    products.push(product);
  }

  return products;
}

// Generate Addresses
async function createAddresses(
  prisma: TxPrisma,
  userAccounts: UserAccount[],
  count: number
) {
  const addresses: Address[] = [];

  for (const userAccount of userAccounts) {
    const addressCount = Math.floor(Math.random() * count) + 1;
    for (let i = 0; i < addressCount; i++) {
      const address = await prisma.address.create({
        data: {
          user_id: userAccount.id,
          full_name: `User ${userAccount.id}`,
          phone: userAccount.phone,
          address: `${Math.floor(Math.random() * 100)} Street`,
          city: "Ho Chi Minh City",
          province: "Ho Chi Minh",
          country: "Vietnam",
        },
      });
      addresses.push(address);
    }
  }

  return addresses;
}

// Generate Carts and Items
async function createCarts(
  prisma: TxPrisma,
  userAccounts: UserAccount[],
  products: Product[]
) {
  const cartsData: Prisma.CartCreateManyInput[] = [];
  const itemOnCartData: Prisma.ItemOnCartCreateManyInput[] = [];

  for (const userAccount of userAccounts) {
    cartsData.push({
      id: userAccount.id,
    });

    // Create cart items for each user
    const itemCount = Math.floor(Math.random() * 5) + 1;
    for (let i = 0; i < itemCount; i++) {
      itemOnCartData.push({
        cart_id: userAccount.id,
        product_id: products[Math.floor(Math.random() * products.length)].id,
        quantity: Math.floor(Math.random() * 5) + 1,
      });
    }
  }

  await prisma.cart.createMany({
    data: cartsData,
    skipDuplicates: true,
  });

  if (itemOnCartData.length > 0) {
    await prisma.itemOnCart.createMany({
      data: itemOnCartData,
      skipDuplicates: true,
    });
  }

  return await prisma.cart.findMany({ take: userAccounts.length });
}

// Generate Sales
async function createSales(
  prisma: TxPrisma,
  productModels: ProductModel[],
  tags: Tag[],
  brands: Brand[],
  count: number
) {
  const sales: Sale[] = [];

  for (let i = 0; i < count; i++) {
    const saleType = faker.helpers.enumValue(SaleType);
    let itemId: bigint;

    // Set item_id based on sale type
    switch (saleType) {
      case SaleType.PRODUCT_MODEL:
        itemId =
          productModels[Math.floor(Math.random() * productModels.length)].id;
        break;
      case SaleType.BRAND:
        itemId = brands[Math.floor(Math.random() * brands.length)].id;
        break;
      case SaleType.TAG:
        itemId = BigInt(tags[Math.floor(Math.random() * tags.length)].id);
        break;
      default:
        throw new Error(`Invalid sale type: ${saleType}`);
    }

    const sale = await prisma.sale.create({
      data: {
        type: saleType,
        item_id: itemId,
        date_started: new Date(),
        date_ended:
          Math.random() > 0.5
            ? new Date(Date.now() + 7 * 24 * 60 * 60 * 1000)
            : null, // 50% chance of having end date
        is_active: true,
        discount_percent:
          Math.random() > 0.5 ? Math.floor(Math.random() * 50) + 1 : null, // 50% chance of percent discount
        discount_price:
          Math.random() > 0.5
            ? BigInt(Math.floor(Math.random() * 1000000))
            : null, // 50% chance of fixed discount
        max_discount_price: BigInt(Math.floor(Math.random() * 2000000)),
        SaleTracking: {
          create: {
            current_stock: BigInt(Math.floor(Math.random() * 100) + 1),
            used: BigInt(0),
          },
        },
      },
      include: {
        SaleTracking: true,
      },
    });

    sales.push(sale);
  }

  return sales;
}

// Generate Payments and ProductOnPayment
async function createPayments(
  prisma: TxPrisma,
  userAccounts: UserAccount[],
  products: Product[],
  count: number
) {
  const payments: Payment[] = [];

  for (let i = 0; i < count; i++) {
    const userAccount =
      userAccounts[Math.floor(Math.random() * userAccounts.length)];
    const productCount = Math.floor(Math.random() * 5) + 1;
    const selectedProducts = products
      .sort(() => Math.random() - 0.5)
      .slice(0, productCount);

    const payment = await prisma.payment.create({
      data: {
        user_id: userAccount.id,
        method: faker.helpers.enumValue(PaymentMethod),
        status: faker.helpers.enumValue(Status),
        address: `Address ${i + 1}`,
        total: BigInt(Math.floor(Math.random() * 1000000)),
        products: {
          create: selectedProducts.map((product) => ({
            product_id: product.id,
            quantity: Math.floor(Math.random() * 5) + 1,
            price: BigInt(Math.floor(Math.random() * 100000)),
            total_price: BigInt(Math.floor(Math.random() * 1000000)),
          })),
        },
      },
    });

    payments.push(payment);
  }

  return payments;
}

// Generate Refunds
async function createRefunds(
  prisma: TxPrisma,
  payments: Payment[],
  adminAccounts: AdminAccount[],
  count: number
) {
  const refunds: Refund[] = [];

  for (let i = 0; i < count; i++) {
    const payment = payments[Math.floor(Math.random() * payments.length)];
    const productOnPayment = await prisma.productOnPayment.findFirst({
      where: {
        payment_id: payment.id,
      },
      orderBy: {
        id: "desc",
      },
    });
    const adminAccount =
      adminAccounts[Math.floor(Math.random() * adminAccounts.length)];

    if (productOnPayment) {
      try {
        const refund = await prisma.refund.create({
          data: {
            product_on_payment_id: productOnPayment.id,
            method: faker.helpers.enumValue(RefundMethod),
            status: faker.helpers.enumValue(Status),
            reason: `Refund reason ${i + 1}`,
            address: `Refund address ${i + 1}`,
            amount: BigInt(Math.floor(Math.random() * 1000000)),
            approved_by_id: adminAccount.id,
          },
        });
        refunds.push(refund);
      } catch (error) {
        // skip duplicate refund
      }
    }
  }

  return refunds;
}

// Generate Comments and Resources
async function createComments(
  prisma: TxPrisma,
  userAccounts: UserAccount[],
  count: number
) {
  const commentsData: Prisma.CommentCreateManyInput[] = [];
  const resourcesData: Prisma.ResourceCreateManyInput[] = [];
  const products = await prisma.product.findMany();
  const productModels = await prisma.productModel.findMany();
  const brands = await prisma.brand.findMany();

  const addComment = (commentId: bigint, destId: bigint, type: CommentType) => {
    const account =
      userAccounts[Math.floor(Math.random() * userAccounts.length)];
    commentsData.push({
      id: commentId,
      type: type,
      account_id: account.id,
      dest_id: destId,
      body: faker.lorem.paragraph(),
      upvote: BigInt(Math.floor(Math.random() * 100)),
      downvote: BigInt(Math.floor(Math.random() * 50)),
      score: Math.floor(Math.random() * 100) - 50,
      date_created: faker.date.past(),
      date_updated: faker.date.recent(),
    });

    // Add 0-3 resources per comment
    const resourceCount = Math.floor(Math.random() * 4);
    for (let j = 0; j < resourceCount; j++) {
      resourcesData.push({
        type: ResourceType.COMMENT,
        owner_id: commentId,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
        order: j + 1,
      });
    }
  };

  let commentId = BigInt(1);

  for (let i = 0; i < count; i++) {
    addComment(
      commentId++,
      productModels[Math.floor(Math.random() * productModels.length)].id,
      CommentType.PRODUCT_MODEL
    );
    addComment(
      commentId++,
      brands[Math.floor(Math.random() * brands.length)].id,
      CommentType.BRAND
    );
    addComment(
      commentId++,
      products[Math.floor(Math.random() * products.length)].id,
      CommentType.COMMENT
    );
  }

  // Set the next value for the comment sequence
  await prisma.$executeRawUnsafe(
    `ALTER SEQUENCE product.comment_id_seq RESTART WITH ${commentId}`
  );

  await prisma.comment.createMany({
    data: commentsData,
    skipDuplicates: true,
  });

  if (resourcesData.length > 0) {
    await prisma.resource.createMany({
      data: resourcesData,
      skipDuplicates: true,
    });
  }

  return await prisma.comment.findMany({ take: count });
}

// Generate Resources for brands, productModels, products, and refunds
async function createResources(
  prisma: TxPrisma,
  brands: Brand[],
  productModels: ProductModel[],
  products: Product[],
  refunds: Refund[],
  count: number
) {
  const resourcesData: Prisma.ResourceCreateManyInput[] = [];

  // Add resources for brands
  for (const brand of brands) {
    // Add 1-3 resources per brand
    const resourceCount = Math.floor(Math.random() * 3) + 1;
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        type: ResourceType.BRAND,
        owner_id: brand.id,
        url: faker.image.url({ width: 800, height: 600 }),
        order: i + 1,
      });
    }
  }

  // Add resources for product models
  for (const productModel of productModels) {
    // Add 2-5 resources per product model
    const resourceCount = Math.floor(Math.random() * 4) + 2;
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        type: ResourceType.PRODUCT_MODEL,
        owner_id: productModel.id,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
        order: i + 1,
      });
    }
  }

  // Add resources for products
  for (const product of products) {
    // Add 1-3 resources per product
    const resourceCount = Math.floor(Math.random() * 3) + 1;
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        type: ResourceType.PRODUCT,
        owner_id: product.id,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
        order: i + 1,
      });
    }
  }

  // Add resources for refunds
  for (const refund of refunds) {
    // Add 0-2 resources per refund (some refunds might not have resources)
    const resourceCount = Math.floor(Math.random() * 3);
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        type: ResourceType.REFUND,
        owner_id: refund.id,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
        order: i + 1,
      });
    }
  }

  await prisma.resource.createMany({
    data: resourcesData,
    skipDuplicates: true,
  });

  // Set the next value for the resource sequence
  await prisma.$executeRawUnsafe(
    `ALTER SEQUENCE product.resource_id_seq RESTART WITH ${BigInt(
      resourcesData.length + 1
    )}`
  );

  return await prisma.resource.findMany({ take: resourcesData.length });
}

// Define preset types and configurations
type PresetType = "light" | "medium" | "heavy";

interface SeedConfig {
  accounts: number;
  brands: number;
  tags: number;
  productTypes: number;
  productModels: number;
  products: number;
  addresses: number;
  sales: number;
  payments: number;
  comments: number;
}

// Update main function to accept preset parameter
async function main(preset: PresetType = "light") {
  const prisma = new PrismaClient();
  const config = SEED_PRESETS[preset];

  try {
    console.log(`Starting to seed database with ${preset} preset...`);

    await prisma.$transaction(
      async (tx) => {
        // Create base data
        const roles = await createRoles(tx);
        const { adminAccounts, userAccounts } = await createAccounts(
          tx,
          config.accounts
        );
        const brands = await createBrands(tx, config.brands);
        const tags = await createTags(tx, config.tags);
        const productTypes = await createProductTypes(tx, config.productTypes);
        const productModels = await createProductModels(
          tx,
          brands,
          tags,
          productTypes,
          config.productModels
        );
        const products = await createProducts(
          tx,
          productModels,
          config.products
        );

        // Create related data
        const addresses = await createAddresses(
          tx,
          userAccounts,
          config.addresses
        );
        const carts = await createCarts(tx, userAccounts, products);
        const sales = await createSales(
          tx,
          productModels,
          tags,
          brands,
          config.sales
        );
        const payments = await createPayments(
          tx,
          userAccounts,
          products,
          config.payments
        );
        const refunds = await createRefunds(
          tx,
          payments,
          adminAccounts,
          Math.floor(config.payments * 0.1)
        ); // 10% of payments
        const comments = await createComments(
          tx,
          userAccounts,
          config.comments
        );

        // Create resources
        const resources = await createResources(
          tx,
          brands,
          productModels,
          products,
          refunds,
          Math.floor(config.productModels * 2) // Average 2 resources per product model
        );

        console.log(`Seeding completed successfully with ${preset} preset!`);
      },
      {
        timeout: 1000000000,
      }
    );
  } catch (error) {
    console.error("Error seeding database:", error);
  } finally {
    await prisma.$disconnect();
  }
}

const SEED_PRESETS: Record<PresetType, SeedConfig> = {
  light: {
    accounts: 10,
    brands: 5,
    tags: 8,
    productTypes: 5,
    productModels: 20,
    products: 50,
    addresses: 15,
    sales: 10,
    payments: 15,
    comments: 20,
  },
  medium: {
    accounts: 1000, // 1k users
    brands: 50, // 50 different brands
    tags: 100, // 100 different tags
    productTypes: 30, // 30 product types
    productModels: 500, // 500 different product models
    products: 2000, // 2k product variants
    addresses: 1500, // ~1.5 addresses per user
    sales: 200, // 200 different sales/promotions
    payments: 3000, // 3k orders
    comments: 5000, // 5k comments
  },
  heavy: {
    accounts: 10000, // 10k users
    brands: 200, // 200 different brands
    tags: 300, // 300 different tags
    productTypes: 100, // 100 product types
    productModels: 2000, // 2k different product models
    products: 10000, // 10k product variants
    addresses: 15000, // ~1.5 addresses per user
    sales: 1000, // 1k different sales/promotions
    payments: 30000, // 30k orders
    comments: 50000, // 50k comments
  },
};

// You can now call main with different presets
main("light"); // or main('medium') or main('heavy')
