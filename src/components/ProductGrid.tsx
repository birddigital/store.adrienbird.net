import React, { useState, useEffect } from 'react';
import { squarespaceClient } from '../api/squarespace-client';
import { Product, SquarespaceApiError } from '../types/squarespace';

interface ProductGridProps {
  category?: string;
  tag?: string;
  limit?: number;
}

const ProductGrid: React.FC<ProductGridProps> = ({
  category,
  tag,
  limit = 12
}) => {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        setLoading(true);
        setError(null);

        const response = await squarespaceClient.getProducts({
          limit,
          category,
          tag
        });

        setProducts(response.result);
      } catch (err) {
        if (err instanceof SquarespaceApiError) {
          setError(`Failed to load products: ${err.message}`);
        } else {
          setError('An unexpected error occurred while loading products');
        }
        console.error('Product fetch error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, [category, tag, limit]);

  if (loading) {
    return (
      <div className="product-grid loading">
        <div className="loading-spinner">Loading products...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="product-grid error">
        <div className="error-message">{error}</div>
        <button
          onClick={() => window.location.reload()}
          className="retry-button"
        >
          Try Again
        </button>
      </div>
    );
  }

  if (products.length === 0) {
    return (
      <div className="product-grid empty">
        <div className="empty-message">No products found</div>
      </div>
    );
  }

  return (
    <div className="product-grid">
      <div className="products-container">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
    </div>
  );
};

interface ProductCardProps {
  product: Product;
}

const ProductCard: React.FC<ProductCardProps> = ({ product }) => {
  const mainVariant = product.products[0];

  if (!mainVariant) return null;

  const price = mainVariant.pricing.basePrice;
  const salePrice = mainVariant.pricing.salePrice;
  const onSale = mainVariant.pricing.onSale;
  const mainImage = mainVariant.images[0];

  const formatPrice = (amount: string, currency: string) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    }).format(parseFloat(amount));
  };

  return (
    <div className="product-card">
      <div className="product-image">
        {mainImage ? (
          <img
            src={mainImage.url}
            alt={mainImage.description || mainVariant.name}
            width={mainImage.width}
            height={mainImage.height}
            loading="lazy"
          />
        ) : (
          <div className="placeholder-image">No Image</div>
        )}

        {onSale && (
          <div className="sale-badge">Sale</div>
        )}
      </div>

      <div className="product-info">
        <h3 className="product-name">{mainVariant.name}</h3>

        {mainVariant.description && (
          <p className="product-description">{mainVariant.description}</p>
        )}

        <div className="product-pricing">
          {price && (
            <span className="current-price">
              {formatPrice(price.value, price.currency)}
            </span>
          )}

          {onSale && salePrice && price && (
            <span className="original-price">
              {formatPrice(price.value, price.currency)}
            </span>
          )}
        </div>

        <div className="product-stock">
          {mainVariant.stock.unlimited ? (
            <span className="in-stock">In Stock</span>
          ) : mainVariant.stock.quantity && mainVariant.stock.quantity > 0 ? (
            <span className="in-stock">
              {mainVariant.stock.quantity} in stock
            </span>
          ) : (
            <span className="out-of-stock">Out of Stock</span>
          )}
        </div>

        <button
          className="add-to-cart-button"
          disabled={!mainVariant.stock.unlimited && mainVariant.stock.quantity === 0}
        >
          {mainVariant.stock.unlimited ||
           (mainVariant.stock.quantity && mainVariant.stock.quantity > 0)
            ? 'Add to Cart'
            : 'Out of Stock'
          }
        </button>
      </div>
    </div>
  );
};

export default ProductGrid;